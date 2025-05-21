package crawler

import (
	"math/rand"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/net/publicsuffix"
)

// Common browser user agents for rotation
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 OPR/82.0.4227.44",
}

// Common accept-language headers for rotation
var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-GB,en;q=0.9",
	"en-CA,en;q=0.9,fr-CA;q=0.8",
	"en;q=0.9",
	"en-US,en;q=0.8,de;q=0.5",
}

// Common accept headers for rotation
var acceptHeaders = []string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
}

// Common accept-encoding headers
var acceptEncodings = []string{
	"gzip, deflate, br",
	"gzip, deflate",
	"br;q=1.0, gzip;q=0.8, *;q=0.1",
}

// AntiBotConfig contains configuration for anti-bot measures
type AntiBotConfig struct {
	EnableUserAgentRotation   bool
	EnableHeaderRandomization bool
	EnableCookieHandling      bool
	EnableSophisticatedDelays bool
	CustomUserAgents          []string
	CustomAcceptLanguages     []string
	BaseDelay                 time.Duration
	RandomDelayFactor         float64 // 0.5-2.0 recommended
}

// NewDefaultAntiBotConfig creates a default anti-bot configuration
func NewDefaultAntiBotConfig() *AntiBotConfig {
	return &AntiBotConfig{
		EnableUserAgentRotation:   true,
		EnableHeaderRandomization: true,
		EnableCookieHandling:      true,
		EnableSophisticatedDelays: true,
		BaseDelay:                 200 * time.Millisecond,
		RandomDelayFactor:         1.5,
	}
}

// ApplyAntiBotMeasures applies anti-bot measures to a collector
func ApplyAntiBotMeasures(collector *colly.Collector, config *AntiBotConfig) error {
	// 1. User Agent Rotation
	if config.EnableUserAgentRotation {
		combinedUserAgents := userAgents
		if len(config.CustomUserAgents) > 0 {
			combinedUserAgents = append(combinedUserAgents, config.CustomUserAgents...)
		}

		collector.OnRequest(func(r *colly.Request) {
			userAgent := combinedUserAgents[rand.Intn(len(combinedUserAgents))]
			r.Headers.Set("User-Agent", userAgent)
		})
	}

	// 2. Header Randomization
	if config.EnableHeaderRandomization {
		combinedAcceptLanguages := acceptLanguages
		if len(config.CustomAcceptLanguages) > 0 {
			combinedAcceptLanguages = append(combinedAcceptLanguages, config.CustomAcceptLanguages...)
		}

		collector.OnRequest(func(r *colly.Request) {
			// Set Accept-Language header
			acceptLanguage := combinedAcceptLanguages[rand.Intn(len(combinedAcceptLanguages))]
			r.Headers.Set("Accept-Language", acceptLanguage)

			// Set Accept header
			acceptHeader := acceptHeaders[rand.Intn(len(acceptHeaders))]
			r.Headers.Set("Accept", acceptHeader)

			// Set Accept-Encoding header
			acceptEncoding := acceptEncodings[rand.Intn(len(acceptEncodings))]
			r.Headers.Set("Accept-Encoding", acceptEncoding)

			// Set DNT (Do Not Track) randomly (some browsers send it, some don't)
			if rand.Intn(2) == 0 {
				r.Headers.Set("DNT", "1")
			}

			// Set Sec-Fetch headers (modern browsers)
			if rand.Intn(2) == 0 {
				r.Headers.Set("Sec-Fetch-Dest", "document")
				r.Headers.Set("Sec-Fetch-Mode", "navigate")
				r.Headers.Set("Sec-Fetch-Site", "none")
				r.Headers.Set("Sec-Fetch-User", "?1")
			}
		})
	}

	// 3. Sophisticated Delays (implemented in the LimitRule)
	if config.EnableSophisticatedDelays {
		// This will be handled by the caller with getRandomDelay
	}

	// 4. Cookie Handling
	if config.EnableCookieHandling {
		// Create a cookie jar that handles cookies like a browser
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			return err
		}
		collector.SetCookieJar(jar)
	}

	return nil
}

// GetRandomDelay returns a random delay based on the base delay and a factor
func GetRandomDelay(baseDelay time.Duration, factor float64) time.Duration {
	// Generate a random factor between 0.5*factor and 1.5*factor
	randomFactor := 0.5*factor + rand.Float64()*factor

	// Ensure the factor doesn't go below 0.5
	if randomFactor < 0.5 {
		randomFactor = 0.5
	}

	return time.Duration(float64(baseDelay) * randomFactor)
}

// GetRefererPolicy returns a function that sets a referer based on the previous page
func GetRefererPolicy() func(*colly.Request, *colly.Response) {
	visitedURLs := make(map[string]string)

	return func(req *colly.Request, resp *colly.Response) {
		// If this is a new request (not a redirect or retry)
		if resp != nil {
			visitedURLs[req.URL.String()] = resp.Request.URL.String()
		}

		// For all requests, try to set a referer if we have visited a page before
		for urlPrefix, referer := range visitedURLs {
			if strings.HasPrefix(req.URL.String(), urlPrefix) && referer != "" {
				req.Headers.Set("Referer", referer)
				break
			}
		}
	}
}
