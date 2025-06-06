package crawler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LoadESConfigFromEnv loads Elasticsearch configuration from environment variables
func LoadESConfigFromEnv() ([]string, string, string) {
	elasticsearchAddresses := SplitEnvVar("CRAWLER_ELASTICSEARCH_ADDRESSES", ",")
	elasticsearchUsername := os.Getenv("CRAWLER_ELASTICSEARCH_USERNAME")
	elasticsearchPassword := os.Getenv("CRAWLER_ELASTICSEARCH_PASSWORD")
	return elasticsearchAddresses, elasticsearchUsername, elasticsearchPassword
}

// LoadConfigFromEnv loads crawler configuration from environment variables
func LoadConfigFromEnv() (*CrawlerConfig, error) {
	// Required values
	startURL := os.Getenv("CRAWLER_START_URL")
	if startURL == "" {
		return nil, fmt.Errorf("CRAWLER_START_URL environment variable is required")
	}

	// Optional values with defaults
	var maxDepth *int
	if envMaxDepth := os.Getenv("CRAWLER_MAX_DEPTH"); envMaxDepth != "" {
		if d, err := strconv.Atoi(envMaxDepth); err == nil && d >= 0 {
			maxDepth = &d
		}
	}

	parallelRequests, _ := strconv.Atoi(GetEnvWithDefault("CRAWLER_PARALLEL_REQUESTS", "10"))

	var maxVisits *int
	if envMaxVisits := os.Getenv("CRAWLER_MAX_VISITS"); envMaxVisits != "" {
		if v, err := strconv.Atoi(envMaxVisits); err == nil && v > 0 {
			maxVisits = &v
		}
	}

	respectRobotsTxt, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_RESPECT_ROBOTS_TXT", "true"))
	delayMs, _ := strconv.Atoi(GetEnvWithDefault("CRAWLER_DELAY_MS", "50"))
	randomDelayMs, _ := strconv.Atoi(GetEnvWithDefault("CRAWLER_RANDOM_DELAY_MS", "50"))
	timeoutSec, _ := strconv.Atoi(GetEnvWithDefault("CRAWLER_TIMEOUT_SEC", "10"))
	ignoreQueryStrings, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_IGNORE_QUERY_STRINGS", "false"))

	// Lists
	allowedDomains := SplitEnvVar("CRAWLER_ALLOWED_DOMAINS", ",")
	disallowedDomains := SplitEnvVar("CRAWLER_DISALLOWED_DOMAINS", ",")
	allowedURLs := SplitEnvVar("CRAWLER_ALLOWED_URLS", ",")
	disallowedURLs := SplitEnvVar("CRAWLER_DISALLOWED_URLS", ",")

	// Proxies
	proxies := SplitEnvVar("CRAWLER_PROXIES", ",")

	// User agent
	userAgent := GetEnvWithDefault("CRAWLER_USER_AGENT",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// Accept-Language header
	acceptLanguage := GetEnvWithDefault("CRAWLER_ACCEPT_LANGUAGE", "")

	// Content and metrics configuration
	enableFullContent, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_FULL_CONTENT", "false"))
	enableMetrics, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_METRICS", "false"))

	// Anti-bot configuration
	enableUserAgentRotation, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_USER_AGENT_ROTATION", "true"))
	enableHeaderRandomization, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_HEADER_RANDOMIZATION", "true"))
	enableCookieHandling, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_COOKIE_HANDLING", "true"))
	enableSophisticatedDelays, _ := strconv.ParseBool(GetEnvWithDefault("CRAWLER_ENABLE_SOPHISTICATED_DELAYS", "true"))
	randomDelayFactor, _ := strconv.ParseFloat(GetEnvWithDefault("CRAWLER_RANDOM_DELAY_FACTOR", "1.5"), 64)
	customUserAgents := SplitEnvVar("CRAWLER_CUSTOM_USER_AGENTS", ",")
	customAcceptLanguages := SplitEnvVar("CRAWLER_CUSTOM_ACCEPT_LANGUAGES", ",")

	// Create anti-bot config
	antiBotConfig := &AntiBotConfig{
		EnableUserAgentRotation:   enableUserAgentRotation,
		EnableHeaderRandomization: enableHeaderRandomization,
		EnableCookieHandling:      enableCookieHandling,
		EnableSophisticatedDelays: enableSophisticatedDelays,
		CustomUserAgents:          customUserAgents,
		CustomAcceptLanguages:     customAcceptLanguages,
		BaseDelay:                 time.Duration(delayMs) * time.Millisecond,
		RandomDelayFactor:         randomDelayFactor,
	}

	return &CrawlerConfig{
		AcceptLanguage:     acceptLanguage,
		StartURL:           startURL,
		MaxDepth:           maxDepth,
		UserAgent:          userAgent,
		ParallelRequests:   parallelRequests,
		MaxVisits:          maxVisits,
		RespectRobotsTxt:   respectRobotsTxt,
		Delay:              time.Duration(delayMs) * time.Millisecond,
		RandomDelay:        time.Duration(randomDelayMs) * time.Millisecond,
		Timeout:            time.Duration(timeoutSec) * time.Second,
		IgnoreQueryStrings: ignoreQueryStrings,
		AllowedDomains:     allowedDomains,
		DisallowedDomains:  disallowedDomains,
		AllowedURLs:        allowedURLs,
		DisallowedURLs:     disallowedURLs,
		Proxies:            proxies,
		EnableFullContent:  enableFullContent,
		EnableMetrics:      enableMetrics,
		AntiBotConfig:      antiBotConfig,
	}, nil
}

// Helper function to get env var with default value
func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to split env var into a slice
func SplitEnvVar(key, separator string) []string {
	value := os.Getenv(key)
	if value == "" {
		return []string{}
	}
	split := strings.Split(value, separator)
	// Trim spaces from each item
	trimmed := make([]string, 0, len(split))
	for _, item := range split {
		item = strings.TrimSpace(item)
		if item != "" {
			trimmed = append(trimmed, item)
		}
	}
	return trimmed
}
