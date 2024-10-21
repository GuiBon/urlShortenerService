package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/config"
	"urlShortenerService/internal/infrastructure/malwarescanner"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"
	"urlShortenerService/internal/transport/http"
	"urlShortenerService/internal/usecase"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Set flag to output glog logs to stderr
	flag.Set("logtostderr", "true")
	flag.Parse()

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

	// Initialize the redis
	statisticsStore, err := statistics.NewRedisStore(cfg.Redis)
	if err != nil {
		log.Fatalf("Error initializing redis: %s", err.Error())
	}

	// Initialize malware scanner
	malwareScanner := malwarescanner.NewDummyScanner()

	// Build the commands
	urlSanitizerCmd := command.URLSanitizerCmdBuilder()
	slugGeneratorCmd := command.SlugGeneratorCmdBuilder(cfg.Slug.MaximalLenght)
	slugValidatorCmd := command.SlugValidatorCmdBuilder(cfg.Slug.MaximalLenght)
	createShortenURLCmd := usecase.CreateShortenURLCmdBuilder(cfg.ServerDomain.CreateBaseURL(), urlSanitizerCmd, slugGeneratorCmd, shortURLStore, statisticsStore)
	getOriginalURLCmd := usecase.GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScanner, shortURLStore, statisticsStore)
	forceGetOriginalURLCmd := usecase.ForceGetOriginalURLCmdBuilder(slugValidatorCmd, shortURLStore, statisticsStore)
	deleteExpiredURLsCmd := usecase.DeleteExpiredURLsCmdBuilder(cfg.Slug.TimeToExpire, shortURLStore)
	getStatisticsForURLCmd := usecase.GetStatisticsForURLCmdBuilder(urlSanitizerCmd, statisticsStore)
	getTopStatisticsCmd := usecase.GetTopStatisticsCmdBuilder(statisticsStore)

	// Build the cron job function
	cronJob := func() {
		glog.Info("cron to delete expired URLs started")
		nbURLsDeleterd, err := deleteExpiredURLsCmd(context.Background())
		if err != nil {
			glog.Warningf("failed to delete expired URLs: %w", err)
		} else {
			glog.Infof("cron to delete expired URLs done. [%d] URLs deleted", nbURLsDeleterd)
		}
	}

	// Run the cron job function at startup
	cronJob()

	// Initialize the cron to delete expired urls and start it
	c := cron.New()
	_, err = c.AddFunc("*/10 * * * *", cronJob) // Every 10 minutes
	if err != nil {
		log.Fatalf("Error initializing delete expired urls cron: %s", err.Error())
	}
	c.Start()

	// Initialize the HTTP router
	router := http.NewBuilder(domain.Environment(os.Getenv("env"))).BuildRouter(createShortenURLCmd, getOriginalURLCmd, forceGetOriginalURLCmd, getStatisticsForURLCmd, getTopStatisticsCmd)

	// Start the service
	router.Run(fmt.Sprintf(":%d", cfg.ServerDomain.Port))
}
