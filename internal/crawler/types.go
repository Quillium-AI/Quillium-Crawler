package crawler

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

// PageData represents the structured data extracted from a webpage
type PageData struct {
	URL     string
	Title   string
	Snippet string
	Images  []string
}

// Crawler represents the web crawler with its collector and queue
type Crawler struct {
	Collector *colly.Collector
	Queue     *queue.Queue
	Domains   []string
	MaxDepth  int
}

// Config holds the configuration for the crawler
type Config struct {
	Domains      []string
	MaxDepth     int
	ThreadCount  int
	MaxQueueSize int
	Parallelism  int
	Delay        time.Duration
	RandomDelay  time.Duration
}
