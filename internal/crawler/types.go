package crawler

import (
	"context"
	"sync"
	"time"

	"github.com/Quillium-AI/Quillium-Crawler/internal/dedup"
	"github.com/gocolly/colly"
)

// Crawler represents a web crawler instance
type Crawler struct {
	Collector *colly.Collector
	Options   *CrawlerConfig
	ctx       context.Context
	cancel    context.CancelFunc
	isRunning bool
	mutex     sync.RWMutex
	wg        sync.WaitGroup
	storage   Storage            // Interface for storage operations
	urlFilter *dedup.BloomFilter // Bloom filter for URL deduplication
}

// Storage defines the interface for storage operations
type Storage interface {
	GetPage(url string) (*PageData, bool)
	SavePage(page *PageData) error
}

// CrawlerConfig contains all configuration options loaded from environment variables
type CrawlerConfig struct {
	AcceptLanguage     string
	StartURL           string
	MaxDepth           *int          // Optional: if nil, no depth limit
	UserAgent          string
	ParallelRequests   int
	MaxVisits          *int          // Optional: if nil, no visit limit
	RespectRobotsTxt   bool
	Delay              time.Duration
	RandomDelay        time.Duration
	Timeout            time.Duration
	IgnoreQueryStrings bool
	AllowedDomains     []string
	DisallowedDomains  []string
	AllowedURLs        []string
	DisallowedURLs     []string
	OutputFile         string
	Proxies            []string
	AntiBotConfig      *AntiBotConfig
	EnableFullContent  bool
	EnableMetrics      bool
}

type CrawlerManager struct {
	crawlers map[string]*Crawler
	mutex    sync.RWMutex
}

// PageData represents the data extracted from a crawled page
type PageData struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Snippet     string `json:"snippet"`
	FullContent string `json:"full_content,omitempty"` // Full HTML content when enabled
}

// JSONStorage handles storing crawled data to a JSON file
type JSONStorage struct {
	filePath string
	mutex    sync.Mutex
}

// NewJSONStorage creates a new JSON storage handler
func NewJSONStorage(filePath string) *JSONStorage {
	return &JSONStorage{
		filePath: filePath,
	}
}

// ProxyManager handles proxy rotation for the crawler
type ProxyManager struct {
	proxies    []string
	currentIdx int
	mutex      sync.Mutex
	enabled    bool
}
