package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

const NSQTopic = "integer_generator"

func main() {
	nsqURL := getRequiredEnv("INTEGER_GEN_NSQ_URL")
	intervalStr := getRequiredEnv("INTEGER_GEN_INTERVAL")

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Error converting \"%s\": %s", intervalStr, err)
	}

	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer(nsqURL, cfg)
	if err != nil {
		log.Fatalf("Cannot connect to NSQd: %s", err)
	}

	go func() {
		tick := time.Tick(interval)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for range tick {
			// Intn returns a number in [0,n) that's why I need to increase 1
			integer := r.Intn(999) + 1
			err := producer.Publish(NSQTopic, []byte{byte(integer)})
			if err != nil {
				log.Println("Cannot publish to NSQd:", err)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("Shutting down,", sig, "signal received")

	producer.Stop()
}

func getRequiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable \"%s\" is empty or missing", name)
	}

	return value
}
