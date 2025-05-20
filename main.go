package main

import (
	"log"

	"github.com/Quillium-AI/Quillium-Crawler/internal/api"
	"github.com/Quillium-AI/Quillium-Crawler/internal/crawler"
	"github.com/gocolly/colly"
)

func main() {
	collector := colly.NewCollector()
	crawler.StartCrawler(collector, "https://quilliumtest.com")
	if err := api.StartServer(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
