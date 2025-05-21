package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

// NewProxyManager creates a new proxy manager
func NewProxyManager(proxies []string) *ProxyManager {
	return &ProxyManager{
		proxies:    proxies,
		currentIdx: 0,
		enabled:    len(proxies) > 0,
	}
}

// ApplyProxy configures the provided colly collector to use the proxy rotation
func (pm *ProxyManager) ApplyProxy(c *colly.Collector) error {
	if !pm.enabled || len(pm.proxies) == 0 {
		return nil // No proxies configured, skip
	}

	// Update the current index for next rotation
	pm.mutex.Lock()
	pm.currentIdx = (pm.currentIdx + 1) % len(pm.proxies)
	pm.mutex.Unlock()

	// Create roundRobin proxy switcher
	proxyRotator, err := proxy.RoundRobinProxySwitcher(pm.proxies...)
	if err != nil {
		return err
	}

	// Set the proxy rotator
	c.SetProxyFunc(proxyRotator)

	return nil
}

// GetCurrentProxy returns the current proxy in the rotation
func (pm *ProxyManager) GetCurrentProxy() string {
	if !pm.enabled || len(pm.proxies) == 0 {
		return ""
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.proxies[pm.currentIdx]
}
