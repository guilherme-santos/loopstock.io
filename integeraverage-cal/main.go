package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	NSQTopic   = "integer_generator"
	APIEnpoint = "/v1/integers"
)

type AverageHandler struct {
	IntegerChan chan int
}

func (h *AverageHandler) HandleMessage(message *nsq.Message) error {
	if len(message.Body) > 1 {
		return errors.New("It was expected a integer but received something else")
	}

	h.IntegerChan <- int(message.Body[0])
	return nil
}

func main() {
	nsqLookupdURL := getRequiredEnv("INTEGERAVERAGE_CAL_NSQLOOKUPD_URL")
	apiURL := getRequiredEnv("INTEGERAVERAGE_CAL_INTEGER_API_URL")
	averageIntervalStr := getRequiredEnv("INTEGERAVERAGE_CAL_AVERAGE_INTERVAL")

	averageInterval, err := time.ParseDuration(averageIntervalStr)
	if err != nil {
		log.Fatalf("Error converting \"%s\": %s", averageIntervalStr, err)
	}

	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(NSQTopic, "average", cfg)
	if err != nil {
		log.Fatalf("Cannot create a NSQ consumer: %s", err)
	}

	integerChan := make(chan int)
	consumer.AddHandler(&AverageHandler{IntegerChan: integerChan})

	err = consumer.ConnectToNSQLookupd(nsqLookupdURL)
	if err != nil {
		log.Fatalf("Cannot connect to NSQd: %s", err)
	}

	apiStopChan := make(chan bool)
	go func() {
		integers := []int{}
		tick := time.Tick(averageInterval)
		for {
			select {
			case <-apiStopChan:
				return
			case i := <-integerChan:
				integers = append(integers, i)
			case <-tick:
				if len(integers) == 0 {
					log.Println("No integer found")
					continue
				}

				average := 0
				for _, v := range integers {
					average += v
				}
				average /= len(integers)

				// clear integer but without destroy previous allocated memory
				integers = integers[:0]

				log.Println("Average is:", average)
				integer := map[string]interface{}{"integer": average}
				body, err := json.Marshal(&integer)
				if err != nil {
					log.Println("Cannot convert body:", err)
					continue
				}

				resp, err := http.Post(apiURL+APIEnpoint, "application/json", bytes.NewReader(body))
				if err != nil {
					log.Printf("Cannot POST to \"%s\": %s", apiURL+APIEnpoint, err)
					continue
				}

				if resp.StatusCode != http.StatusCreated {
					bodyResp, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Cannot read API response: %s", err)
					}

					log.Printf("Error adding integer status[%s]: %s", resp.Status, bodyResp)
				}

				resp.Body.Close()
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-consumer.StopChan:
			return
		case sig := <-sigChan:
			log.Println("Shutting down,", sig, "signal received")
			apiStopChan <- true
			consumer.Stop()
		}
	}
}

func getRequiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable \"%s\" is empty or missing", name)
	}

	return value
}
