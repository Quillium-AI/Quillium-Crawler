package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Quillium-AI/Quillium-Crawler/internal/api"
	"github.com/Quillium-AI/Quillium-Crawler/internal/crawler"
	"github.com/Quillium-AI/Quillium-Crawler/internal/elasticsearch"
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

	// Initialize ElasticSearch configuration
	elasticAddresses, _, _ := crawler.LoadESConfigFromEnv()
	log.Printf("ElasticSearch configuration loaded. Addresses: %v", elasticAddresses)
}

func main() {
	// Initialize ElasticSearch client
	elasticAddresses, elasticUsername, elasticPassword := crawler.LoadESConfigFromEnv()
	elasticClient, err := elasticsearch.Initialize(elasticAddresses, elasticUsername, elasticPassword)
	if err != nil {
		log.Fatalf("Failed to initialize ElasticSearch client: %v", err)
	}

	// Get index name from config
	indexName := crawler.GetEnvWithDefault("CRAWLER_INDEX_NAME", "crawled_data")
	log.Printf("Using ElasticSearch index: %s", indexName)

	// Wait for Elasticsearch to be ready with retry mechanism
	log.Println("Waiting for Elasticsearch to be ready...")
	maxRetries := 10
	initialBackoff := 2 * time.Second
	if err := elasticsearch.WaitForElasticsearch(elasticClient, maxRetries, initialBackoff); err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Create ElasticSearch storage
	esStorage := elasticsearch.NewESStorage(elasticClient, indexName)

	// Initialize the index if it doesn't exist
	log.Println("Initializing Elasticsearch index...")
	if err := esStorage.InitializeIndex(); err != nil {
		log.Printf("Warning: Failed to initialize Elasticsearch index: %v", err)
	}

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
		cfgCopy.IndexName = indexName
		crawlerID := "crawler_" + strconv.Itoa(i+1)

		// Set the ElasticSearch storage in the crawler config
		cfgCopy.Storage = esStorage

		crawlerInstance := crawler.NewCrawler(&cfgCopy)

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
	if err := api.StartServer(":8090"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
