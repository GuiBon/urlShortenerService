package main

import (
	"fmt"
	"log"
	"os"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/config"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/transport/http"
	"urlShortenerService/internal/usecase"

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

	// Initialize the database
	shortURLStore, err := shorturl.NewPSQLStore(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database [%s]: %s", cfg.Database.DbName, err.Error())
	}

	// Build the commands
	urlSanitizerCmd := command.URLSanitizerCmdBuilder()
	slugGeneratorCmd := command.SlugGeneratorCmdBuilder(cfg.SlugMaximalLenght)
	slugValidatorCmd := command.SlugValidatorCmdBuilder(cfg.SlugMaximalLenght)
	createShortenURLCmd := usecase.CreateShortenURLCmdBuilder(cfg.ServerDomain.CreateBaseURL(), urlSanitizerCmd, slugGeneratorCmd, shortURLStore)
	getOriginalURLCmd := usecase.GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLStore)

	// Initialize the HTTP router
	router := http.NewBuilder(domain.Environment(os.Getenv("env"))).BuildRouter(createShortenURLCmd, getOriginalURLCmd)

	// Start the service
	router.Run(fmt.Sprintf(":%d", cfg.ServerDomain.Port))
}
