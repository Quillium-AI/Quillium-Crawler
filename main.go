package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Quillium-AI/Quillium-Crawler/internal/crawler"
)

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnvOrDefault(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Could not parse %s as integer, using default: %v", key, err)
		return defaultValue
	}
	return value
}

// getEnvAsDuration gets an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnvOrDefault(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: Could not parse %s as duration, using default: %v", key, err)
		return defaultValue
	}
	return value
}

// getEnvAsSlice gets an environment variable as a slice or returns a default value
func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	valueStr := getEnvOrDefault(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, separator)
}

func main() {
	// Load configuration from environment variables
	config := crawler.Config{
		Domains:      getEnvAsSlice("CRAWLER_ALLOWED_DOMAINS", []string{"quilliumtest.com", "quilliumexample.com"}, ","),
		MaxDepth:     getEnvAsInt("CRAWLER_MAX_DEPTH", 3),
		ThreadCount:  getEnvAsInt("CRAWLER_THREAD_COUNT", 2),
		MaxQueueSize: getEnvAsInt("CRAWLER_MAX_QUEUE_SIZE", 10000),
		Parallelism:  getEnvAsInt("CRAWLER_PARALLELISM", 2),
		Delay:        getEnvAsDuration("CRAWLER_DELAY", 1*time.Second),
		RandomDelay:  getEnvAsDuration("CRAWLER_RANDOM_DELAY", 1*time.Second),
	}

	// Create a new crawler
	crawler, err := crawler.NewCrawler(config)
	if err != nil {
		log.Fatalf("Failed to create crawler: %v", err)
	}

	// Setup handlers for processing web pages
	crawler.SetupHandlers()

	// Add starting URLs to the queue
	startURLs := getEnvAsSlice("CRAWLER_START_URLS", []string{
		"https://quilliumtest.com",
		"https://quilliumexample.com",
	}, ",")

	for _, url := range startURLs {
		if err := crawler.AddURL(url); err != nil {
			log.Printf("Failed to add URL %s: %v", url, err)
		}
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the crawler in a goroutine
	go func() {
		fmt.Println("Starting Quillium Crawler...")
		if err := crawler.Run(config.Parallelism, config.Delay, config.RandomDelay); err != nil {
			log.Printf("Crawler error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down Quillium Crawler...")
}