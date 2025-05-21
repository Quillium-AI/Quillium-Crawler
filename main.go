package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Quillium-AI/Quillium-Crawler/internal/api"
	"github.com/Quillium-AI/Quillium-Crawler/internal/crawler"
)

// Global configuration
var config *crawler.CrawlerConfig

// init loads configuration from environment variables
func init() {
	var err error
	config, err = crawler.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded. Starting URL: %s, Max Depth: %d",
		config.StartURL, config.MaxDepth)

	// Initialize the storage file
	storage := crawler.NewJSONStorage(config.OutputFile)
	if err := storage.Initialize(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
}

func main() {
	// Create crawler manager
	manager := crawler.NewCrawlerManager()

	// Support multiple start URLs (comma-separated)
	startURLs := []string{config.StartURL}
	if envURLs := os.Getenv("CRAWLER_START_URLS"); envURLs != "" {
		startURLs = crawler.SplitEnvVar("CRAWLER_START_URLS", ",")
	}

	for i, url := range startURLs {
		cfgCopy := *config
		cfgCopy.StartURL = url
		crawlerID := "crawler_" + strconv.Itoa(i+1)

		crawlerInstance := crawler.NewCrawler(&cfgCopy)

		storage := crawler.NewJSONStorage(cfgCopy.OutputFile)
		storage.RegisterStorageCallbacks(crawlerInstance.Collector)

		if len(cfgCopy.Proxies) > 0 {
			proxyManager := crawler.NewProxyManager(cfgCopy.Proxies)
			if err := proxyManager.ApplyProxy(crawlerInstance.Collector); err != nil {
				log.Printf("Warning: Failed to configure proxies: %v", err)
			} else {
				log.Printf("Configured %d proxies for rotation", len(cfgCopy.Proxies))
			}
		}

		manager.AddCrawler(crawlerID, crawlerInstance)
		manager.StartCrawler(crawlerID)
		log.Printf("Started crawler %s for URL: %s", crawlerID, url)
	}

	// Start API server
	if err := api.StartServer(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
