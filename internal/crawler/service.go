package crawler

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

// DefaultConfig returns the default configuration for the crawler
func DefaultConfig() Config {
	return Config{
		Domains:      []string{"example.com"},
		MaxDepth:     3,
		ThreadCount:  2,
		MaxQueueSize: 10000,
		Parallelism:  2,
		Delay:        1 * time.Second,
		RandomDelay:  1 * time.Second,
	}
}

// NewCrawler creates a new crawler with the specified configuration
func NewCrawler(config Config) (*Crawler, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(config.Domains...),
		colly.MaxDepth(config.MaxDepth),
		colly.Async(true),
		colly.UserAgent(randomUserAgent()),
	)

	q, err := queue.New(
		config.ThreadCount,
		&queue.InMemoryQueueStorage{MaxSize: config.MaxQueueSize},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	return &Crawler{
		Collector: c,
		Queue:     q,
		Domains:   config.Domains,
		MaxDepth:  config.MaxDepth,
	}, nil
}

// SetupHandlers configures the collector's callbacks for processing web pages
func (cr *Crawler) SetupHandlers() {
	cr.Collector.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.DOM.Find("title").Text()
		snippet := e.DOM.Find("p").First().Text()
		images := []string{}
		e.DOM.Find("img").Each(func(_ int, img *goquery.Selection) {
			src, exists := img.Attr("src")
			if exists {
				images = append(images, e.Request.AbsoluteURL(src))
			}
		})

		data := PageData{
			URL:     e.Request.URL.String(),
			Title:   title,
			Snippet: snippet,
			Images:  images,
		}

		cr.processPageData(data)
	})

	cr.Collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", r.Request.URL, err)
	})
}

// AddURL adds a URL to the crawler queue
func (cr *Crawler) AddURL(url string) error {
	return cr.Queue.AddURL(url)
}

// Run starts the crawler and processes the queue
func (cr *Crawler) Run(parallelism int, delay, randomDelay time.Duration) error {
	cr.Collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallelism,
		Delay:       delay,
		RandomDelay: randomDelay,
	})

	return cr.Queue.Run(cr.Collector)
}

// processPageData handles the extracted page data
func (cr *Crawler) processPageData(data PageData) {
	// For now, just print the data
	// In a real implementation, this would send to an indexer or database
	fmt.Printf("Crawled: %s\nTitle: %s\nSnippet: %s\nImages: %d\n\n",
		data.URL, data.Title, data.Snippet, len(data.Images))
}

// Common user agents for the crawler to use
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
}

// randomUserAgent returns a random user agent from the list
func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
