package crawler

import (
	"sync"
)

// NewCrawlerManager creates a new crawler manager
func NewCrawlerManager() *CrawlerManager {
	return &CrawlerManager{
		crawlers: make(map[string]*Crawler),
		mutex:    sync.RWMutex{},
	}
}

// AddCrawler adds a crawler to the manager with the given ID
func (m *CrawlerManager) AddCrawler(id string, crawler *Crawler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.crawlers[id] = crawler
}

// GetCrawler retrieves a crawler by ID
func (m *CrawlerManager) GetCrawler(id string) (*Crawler, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	crawler, exists := m.crawlers[id]
	return crawler, exists
}

// RemoveCrawler removes a crawler from the manager
func (m *CrawlerManager) RemoveCrawler(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if crawler, exists := m.crawlers[id]; exists && crawler.IsRunning() {
		crawler.Stop()
	}
	delete(m.crawlers, id)
}

// StartCrawler starts a crawler by ID
func (m *CrawlerManager) StartCrawler(id string) bool {
	m.mutex.RLock()
	crawler, exists := m.crawlers[id]
	m.mutex.RUnlock()

	if !exists {
		return false
	}

	crawler.Start()
	return true
}

// StopCrawler stops a crawler by ID
func (m *CrawlerManager) StopCrawler(id string) bool {
	m.mutex.RLock()
	crawler, exists := m.crawlers[id]
	m.mutex.RUnlock()

	if !exists {
		return false
	}

	crawler.Stop()
	return true
}

// GetCrawlerStatus returns the status of a crawler
func (m *CrawlerManager) GetCrawlerStatus(id string) (bool, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	crawler, exists := m.crawlers[id]
	if !exists {
		return false, false
	}
	return crawler.IsRunning(), true
}

// GetAllCrawlerIDs returns all crawler IDs
func (m *CrawlerManager) GetAllCrawlerIDs() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	ids := make([]string, 0, len(m.crawlers))
	for id := range m.crawlers {
		ids = append(ids, id)
	}
	return ids
}

// StopAllCrawlers stops all running crawlers
func (m *CrawlerManager) StopAllCrawlers() {
	m.mutex.RLock()
	crawlers := make([]*Crawler, 0, len(m.crawlers))
	for _, crawler := range m.crawlers {
		if crawler.IsRunning() {
			crawlers = append(crawlers, crawler)
		}
	}
	m.mutex.RUnlock()

	for _, crawler := range crawlers {
		crawler.Stop()
	}
}
