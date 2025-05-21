package crawler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

// Initialize creates the JSON file if it doesn't exist and initializes it with an empty array
func (s *JSONStorage) Initialize() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if file exists
	_, err := os.Stat(s.filePath)
	if os.IsNotExist(err) {
		// Create the file with an empty JSON array
		f, err := os.Create(s.filePath)
		if err != nil {
			return fmt.Errorf("failed to create storage file: %v", err)
		}
		defer f.Close()

		// Write empty array
		_, err = f.WriteString("[]")
		if err != nil {
			return fmt.Errorf("failed to initialize storage file: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking storage file: %v", err)
	}

	return nil
}

// StorePageData appends crawled page data to the JSON file
func (s *JSONStorage) StorePageData(data PageData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Read current file content
	content, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read storage file: %v", err)
	}

	// Parse existing data
	var pages []PageData
	err = json.Unmarshal(content, &pages)
	if err != nil {
		return fmt.Errorf("failed to parse storage file: %v", err)
	}

	// Check if page already exists, update if it does
	for i, page := range pages {
		if page.URL == data.URL {
			// Preserve full content if it exists and new data doesn't have it
			if page.FullContent != "" && data.FullContent == "" {
				data.FullContent = page.FullContent
			}
			pages[i] = data
			// Write back to file
			updatedContent, err := json.MarshalIndent(pages, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to serialize data: %v", err)
			}

			err = os.WriteFile(s.filePath, updatedContent, 0644)
			if err != nil {
				return fmt.Errorf("failed to write to storage file: %v", err)
			}

			return nil
		}
	}

	// Append new data if not found
	pages = append(pages, data)

	// Write back to file
	updatedContent, err := json.MarshalIndent(pages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize data: %v", err)
	}

	err = os.WriteFile(s.filePath, updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to storage file: %v", err)
	}

	return nil
}

// RegisterStorageCallbacks registers colly callbacks to store crawled data
func (s *JSONStorage) RegisterStorageCallbacks(c *colly.Collector) {
	c.OnHTML("html", func(e *colly.HTMLElement) {
		title := e.ChildText("title")
		url := e.Request.URL.String()

		// Extract snippet (meta description or first paragraph)
		snippet := e.ChildText("meta[name=description]")
		if snippet == "" {
			snippet = e.ChildText("p")
		}

		// Truncate snippet if too long
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}

		// Store the data
		pageData := PageData{
			URL:     url,
			Title:   title,
			Snippet: snippet,
		}

		err := s.StorePageData(pageData)
		if err != nil {
			fmt.Printf("Error storing page data: %v\n", err)
		}
	})
}

// GetPage retrieves a page by URL from the storage
func (s *JSONStorage) GetPage(url string) (*PageData, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Read current file content
	content, err := os.ReadFile(s.filePath)
	if err != nil {
		fmt.Printf("Error reading storage file: %v\n", err)
		return nil, false
	}

	// Parse existing data
	var pages []PageData
	err = json.Unmarshal(content, &pages)
	if err != nil {
		fmt.Printf("Error parsing storage file: %v\n", err)
		return nil, false
	}

	// Find the page
	for i, page := range pages {
		if page.URL == url {
			return &pages[i], true
		}
	}

	return nil, false
}

// SavePage updates an existing page or adds a new one
func (s *JSONStorage) SavePage(page *PageData) error {
	if page == nil {
		return fmt.Errorf("cannot save nil page")
	}
	return s.StorePageData(*page)
}
