package main

import (
	"log"
	"urlShortenerService/internal/infrastructure/config"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/transport/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	// Initialize the HTTP router
	router := http.NewRouter()

	// Initialize the database
	_, err = shorturl.NewPSQLStore(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database [%s]: %s", cfg.Database.DbName, err.Error())
	}

	// Start the service on the port 8080
	router.Run(":8080")
}
