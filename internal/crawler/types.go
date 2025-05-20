package crawler

import (
	"sync"
	"time"
)

// CrawlerOptions contains configuration options for the crawler
type CrawlerOptions struct {
	StartURL           string
	MaxDepth           int
	UserAgent          string
	ParallelRequests   int
	MaxVisits          int
	RespectRobotsTxt   bool
	Delay              time.Duration
	RandomDelay        time.Duration
	Timeout            time.Duration
	IgnoreQueryStrings bool
	AllowedDomains     []string
	DisallowedDomains  []string
	AllowedURLs        []string
	DisallowedURLs     []string
}

type CrawlerManager struct {
	crawlers map[string]*Crawler
	mutex    sync.RWMutex
}
