package crawler

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// DefaultCrawlerOptions returns a CrawlerOptions with sensible defaults
func DefaultCrawlerOptions() *CrawlerOptions {
	return &CrawlerOptions{
		MaxDepth:           3,
		UserAgent:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		ParallelRequests:   2,
		MaxVisits:          1000,
		RespectRobotsTxt:   true,
		Delay:              200 * time.Millisecond,
		RandomDelay:        100 * time.Millisecond,
		Timeout:            10 * time.Second,
		IgnoreQueryStrings: false,
	}
}

// Crawler represents a web crawler instance
type Crawler struct {
	Collector *colly.Collector
	Options   *CrawlerOptions
	ctx       context.Context
	cancel    context.CancelFunc
	isRunning bool
	mutex     sync.RWMutex
	wg        sync.WaitGroup
}

// NewCrawler creates a new crawler with the given options
func NewCrawler(options *CrawlerOptions) *Crawler {
	if options == nil {
		options = DefaultCrawlerOptions()
	}

	ctx, cancel := context.WithCancel(context.Background())

	collectorOptions := []func(*colly.Collector){
		colly.MaxDepth(options.MaxDepth),
		colly.UserAgent(options.UserAgent),
		colly.Async(true),
	}

	// Handle robots.txt if needed
	if options.RespectRobotsTxt {
		collector := colly.NewCollector()
		collector.AllowURLRevisit = true
		collectorOptions = append(collectorOptions, func(c *colly.Collector) {
			c.AllowURLRevisit = true
		})
	}

	if len(options.AllowedDomains) > 0 {
		collectorOptions = append(collectorOptions, colly.AllowedDomains(options.AllowedDomains...))
	}

	if len(options.DisallowedDomains) > 0 {
		collectorOptions = append(collectorOptions, colly.DisallowedDomains(options.DisallowedDomains...))
	}

	// Handle query string ignoring if needed
	if options.IgnoreQueryStrings {
		collectorOptions = append(collectorOptions, func(c *colly.Collector) {
			c.URLFilters = nil
		})
	}

	collector := colly.NewCollector()
	// Apply all collector options
	for _, option := range collectorOptions {
		option(collector)
	}

	// Set limits
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: options.ParallelRequests,
		Delay:       options.Delay,
		RandomDelay: options.RandomDelay,
	})

	// Add extensions
	extensions.RandomUserAgent(collector)
	extensions.Referer(collector)

	return &Crawler{
		Collector: collector,
		Options:   options,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}
}

// Start begins the crawling process
func (c *Crawler) Start() {
	c.mutex.Lock()
	if c.isRunning {
		c.mutex.Unlock()
		log.Println("Crawler is already running")
		return
	}
	c.isRunning = true
	c.mutex.Unlock()

	c.setupCollector()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() {
			c.mutex.Lock()
			c.isRunning = false
			c.mutex.Unlock()
		}()

		visitCount := 0
		c.Collector.OnRequest(func(r *colly.Request) {
			select {
			case <-c.ctx.Done():
				r.Abort()
				return
			default:
				visitCount++
				if c.Options.MaxVisits > 0 && visitCount > c.Options.MaxVisits {
					r.Abort()
					c.Stop()
					return
				}
				log.Println("Visiting", r.URL.String())
			}
		})

		err := c.Collector.Visit(c.Options.StartURL)
		if err != nil {
			log.Printf("Error starting crawler: %v", err)
		}

		c.Collector.Wait()
	}()
}

// Stop halts the crawling process
func (c *Crawler) Stop() {
	c.mutex.RLock()
	isRunning := c.isRunning
	c.mutex.RUnlock()

	if !isRunning {
		log.Println("Crawler is not running")
		return
	}

	log.Println("Stopping crawler...")
	c.cancel()
	c.wg.Wait()
	log.Println("Crawler stopped")
}

// IsRunning returns whether the crawler is currently active
func (c *Crawler) IsRunning() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isRunning
}

// setupCollector configures the collector with callbacks
func (c *Crawler) setupCollector() {
	c.Collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absURL := e.Request.AbsoluteURL(link)

		// Check if URL is in allowed/disallowed lists
		if len(c.Options.AllowedURLs) > 0 {
			allowed := false
			for _, pattern := range c.Options.AllowedURLs {
				if strings.Contains(absURL, pattern) {
					allowed = true
					break
				}
			}
			if !allowed {
				return
			}
		}

		for _, pattern := range c.Options.DisallowedURLs {
			if strings.Contains(absURL, pattern) {
				return
			}
		}

		log.Printf("Link found: %v -> %v", e.Text, absURL)
		c.Collector.Visit(absURL)
	})

	c.Collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Error visiting %s: %v", r.Request.URL, err)
	})
}

// StartCrawler is a legacy function for backward compatibility
func StartCrawler(collector *colly.Collector, startURL string) {
	options := DefaultCrawlerOptions()
	options.StartURL = startURL

	crawler := NewCrawler(options)
	crawler.Collector = collector // Use the provided collector
	crawler.Start()
}
