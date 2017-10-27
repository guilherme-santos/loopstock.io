package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	dbClient, err := getDatabaseClient()
	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}

	serverPort := getRequiredEnv("INTEGER_API_PORT")

	router := httprouter.New()
	api := NewIntegerAPI(dbClient, router)

	log.Printf("Running webserver on port %s", serverPort)
	err = api.Run(serverPort)
	if err != nil {
		log.Fatal("Cannot run webserver: %s", err)
	}
}

func getRequiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable \"%s\" is empty or missing", name)
	}

	return value
}

func getDatabaseClient() (*sql.DB, error) {
	dbHost := getRequiredEnv("INTEGER_API_DATABASE_HOST")
	dbPort := getRequiredEnv("INTEGER_API_DATABASE_PORT")
	dbUser := getRequiredEnv("INTEGER_API_DATABASE_USER")
	dbPassword := getRequiredEnv("INTEGER_API_DATABASE_PASSWORD")
	dbDatabase := getRequiredEnv("INTEGER_API_DATABASE_DB")

	// hiding password
	log.Printf("Trying to connect to MySQL in '%s@%s:%s/%s'...", dbUser, dbHost, dbPort, dbDatabase)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Cannot ping MySQL: %s", err.Error())
	}

	return db, nil
}
