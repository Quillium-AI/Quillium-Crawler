package crawler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.cherkaoui.ch/quillium-ai/quillium-crawler/internal/dedup"
	"gitlab.cherkaoui.ch/quillium-ai/quillium-crawler/internal/elasticsearch"
	"gitlab.cherkaoui.ch/quillium-ai/quillium-crawler/internal/metrics"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// NewCrawler creates a new crawler with the given options
func NewCrawler(options *CrawlerConfig) *Crawler {

	ctx, cancel := context.WithCancel(context.Background())

	// Build all collector options in one place to avoid duplication
	collectorOptions := []func(*colly.Collector){
		colly.UserAgent(options.UserAgent),
		colly.Async(true), // Enable async for better performance
	}

	// Only set MaxDepth if it's explicitly set
	if options.MaxDepth != nil {
		collectorOptions = append(collectorOptions, colly.MaxDepth(*options.MaxDepth))
	}

	// By default Colly respects robots.txt unless IgnoreRobotsTxt is explicitly set
	if !options.RespectRobotsTxt {
		collectorOptions = append(collectorOptions, colly.IgnoreRobotsTxt())
	} else {
		// Allow URL revisit when respecting robots.txt
		collectorOptions = append(collectorOptions, func(c *colly.Collector) {
			c.AllowURLRevisit = true
		})
	}

	// Add domain filters
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

	// Create the collector with all options at once for better initialization
	collector := colly.NewCollector(collectorOptions...)

	// Enable DNS caching to reduce network latency
	collector.WithTransport(&http.Transport{
		DisableKeepAlives:   false, // Enable keep-alives for connection reuse
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	})

	// Set optimized limits with per-domain parallelism for better performance
	// This allows more concurrent requests while still respecting per-domain rate limits
	if options.AntiBotConfig != nil && options.AntiBotConfig.EnableSophisticatedDelays {
		// Global limit rule for all domains
		collector.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: options.ParallelRequests,
			Delay:       options.Delay,
			RandomDelay: GetRandomDelay(options.AntiBotConfig.BaseDelay, options.AntiBotConfig.RandomDelayFactor),
		})

		// Add specific rules for common domains to allow more parallel requests
		// while still respecting the global delay settings
		if len(options.AllowedDomains) > 0 {
			for _, domain := range options.AllowedDomains {
				collector.Limit(&colly.LimitRule{
					DomainGlob:  domain,
					Parallelism: options.ParallelRequests * 2, // Double parallelism for known domains
					Delay:       options.Delay,
					RandomDelay: GetRandomDelay(options.AntiBotConfig.BaseDelay, options.AntiBotConfig.RandomDelayFactor),
				})
			}
		}
	} else {
		// Use standard delay settings with optimized parallelism
		collector.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: options.ParallelRequests,
			Delay:       options.Delay,
			RandomDelay: options.RandomDelay,
		})
	}

	// Apply anti-bot measures if configured
	if options.AntiBotConfig != nil {
		if err := ApplyAntiBotMeasures(collector, options.AntiBotConfig); err != nil {
			log.Printf("Warning: Failed to apply anti-bot measures: %v", err)
		}
	} else {
		// Add basic extensions if anti-bot is not configured
		extensions.RandomUserAgent(collector)
		extensions.Referer(collector)

		// Set Accept-Language header if configured
		if options.AcceptLanguage != "" {
			collector.OnRequest(func(r *colly.Request) {
				r.Headers.Set("Accept-Language", options.AcceptLanguage)
			})
		}
	}

	// Initialize metrics if enabled
	if options.EnableMetrics {
		metrics.Initialize(options.EnableFullContent)
	}

	// Calculate optimal bloom filter size based on expected visits
	// Use a larger size for more accuracy, smaller size for less memory usage
	expectedURLs := 10000 // Default minimum size to avoid too many false positives
	if options.MaxVisits != nil && *options.MaxVisits > 0 {
		expectedURLs = *options.MaxVisits * 10 // Estimate 10x the max visits as potential URLs
	}

	// Create bloom filter with 1% false positive rate
	bloomSize := dedup.CalculateOptimalSize(expectedURLs, 0.01)
	hashFuncs := dedup.CalculateOptimalHashFunctions(bloomSize, expectedURLs)

	// Debug: log.Printf("Initialized URL filter with size %d and %d hash functions", bloomSize, hashFuncs)

	return &Crawler{
		Collector: collector,
		Options:   options,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
		Storage:   options.Storage,
		urlFilter: dedup.NewBloomFilter(bloomSize, hashFuncs),
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
				if c.Options.MaxVisits != nil && visitCount > *c.Options.MaxVisits {
					r.Abort()
					c.Stop()
					return
				}
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
	// Extract and store page data for all pages
	c.Collector.OnHTML("html", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		// Extract page title
		title := e.ChildText("title")
		if title == "" {
			title = e.ChildText("h1")
		}

		// Create a snippet from meta description or first paragraph
		snippet := e.ChildAttr("meta[name=description]", "content")
		if snippet == "" {
			snippet = e.ChildText("p")
		}
		// Trim snippet if it's too long
		if len(snippet) > 500 {
			snippet = snippet[:500] + "..."
		}

		// Prepare page data
		pageData := &elasticsearch.PageData{
			URL:     url,
			Title:   title,
			Snippet: snippet,
		}

		// Add full content if enabled
		if c.Options.EnableFullContent {
			pageData.FullContent = string(e.Response.Body)

			// Track content size for metrics
			if c.Options.EnableMetrics {
				metrics.ContentSize.Observe(float64(len(e.Response.Body)))
			}
		}

		// Save to storage
		if err := c.Storage.SavePage(pageData); err != nil {
			log.Printf("Error saving page data for %s: %v", url, err)
		}

		// Increment metrics counter
		if c.Options.EnableMetrics {
			metrics.PagesCrawled.Inc()
		}
	})

	// Error handling with more context
	c.Collector.OnError(func(r *colly.Response, err error) {
		if r == nil {
			log.Printf("Request failed: %v", err)
			return
		}

		status := "unknown"
		if r.StatusCode > 0 {
			status = strconv.Itoa(r.StatusCode)
		}

		log.Printf("Request failed for %s (Status: %s): %v",
			r.Request.URL, status, err)

		// Increment error counter if metrics are enabled
		if c.Options.EnableMetrics {
			metrics.RequestErrors.Inc()
		}
	})

	// Success callback for metrics
	c.Collector.OnResponse(func(r *colly.Response) {
		if c.Options.EnableMetrics {
			metrics.RequestsTotal.Inc()
			if r.StatusCode > 0 {
				metrics.RequestsByStatus.WithLabelValues(strconv.Itoa(r.StatusCode)).Inc()
			}
		}
	})

	c.Collector.OnHTML("a[href]", c.handleLink)
	c.Collector.OnError(c.handleError)
}

// handleLink processes discovered links and applies filtering before visiting
func (c *Crawler) handleLink(e *colly.HTMLElement) {
	link := e.Attr("href")
	absURL := e.Request.AbsoluteURL(link)

	// Skip if URL doesn't match allowed patterns
	if !c.isAllowedURL(absURL) {
		return
	}

	// Skip URLs we've already seen (deduplication using bloom filter)
	if c.urlFilter.Contains(absURL) {
		// Debug: log.Printf("Skipping already visited URL: %s", absURL)
		return
	}

	// Mark URL as seen before visiting
	c.urlFilter.Add(absURL)

	// Debug: log.Printf("Link found: %v -> %v", e.Text, absURL)
	c.Collector.Visit(absURL)
}

// isAllowedURL checks allowed/disallowed URL patterns
func (c *Crawler) isAllowedURL(url string) bool {
	if len(c.Options.AllowedURLs) > 0 {
		allowed := false
		for _, pattern := range c.Options.AllowedURLs {
			if strings.Contains(url, pattern) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	for _, pattern := range c.Options.DisallowedURLs {
		if strings.Contains(url, pattern) {
			return false
		}
	}
	return true
}

// handleError logs crawling errors
func (c *Crawler) handleError(r *colly.Response, err error) {
	log.Printf("Error visiting %s: %v", r.Request.URL, err)
}
