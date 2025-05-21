package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics for the crawler
var (
	// RequestsTotal counts the total number of requests
	RequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "crawler_requests_total",
			Help: "Total number of HTTP requests made",
		},
	)

	// RequestsByStatus counts requests by status code
	RequestsByStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "crawler_requests_by_status_total",
			Help: "Number of requests by status code",
		},
		[]string{"status"}, // Only status code as label
	)

	// RequestErrors counts the number of failed requests
	RequestErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "crawler_request_errors_total",
			Help: "Total number of failed requests",
		},
	)

	// PagesCrawled counts the number of pages crawled
	PagesCrawled = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "crawler_pages_crawled_total",
			Help: "Total number of pages successfully crawled",
		},
	)

	// ContentSize tracks the size of crawled content
	ContentSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "crawler_content_size_bytes",
			Help:    "Size of crawled content in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6), // 100B to 1GB
		},
	)

	// FullContentEnabled tracks if full content scraping is enabled
	FullContentEnabled = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "crawler_full_content_enabled",
			Help: "Indicates if full content scraping is enabled (1) or disabled (0)",
		},
	)
)

// Initialize sets up the initial state of metrics
func Initialize(enableFullContent bool) {
	// Set initial values
	if enableFullContent {
		FullContentEnabled.Set(1)
	} else {
		FullContentEnabled.Set(0)
	}
}
