package elasticsearch

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// Initialize creates the Elasticsearch client
func Initialize(addresses []string, username, password string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}
	return client, nil
}

// WaitForElasticsearch waits for Elasticsearch to be ready with exponential backoff
func WaitForElasticsearch(client *elasticsearch.Client, maxRetries int, initialBackoff time.Duration) error {
	var lastError error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		// Check if Elasticsearch is ready by pinging it
		res, err := client.Ping()
		if err == nil && !res.IsError() {
			log.Printf("Successfully connected to Elasticsearch after %d attempts", i+1)
			return nil
		}

		if err != nil {
			lastError = err
			log.Printf("Elasticsearch not ready (attempt %d/%d): %v", i+1, maxRetries, err)
		} else {
			lastError = fmt.Errorf("elasticsearch returned error status: %s", res.String())
			log.Printf("Elasticsearch not ready (attempt %d/%d): %s", i+1, maxRetries, res.String())
		}

		// Sleep with exponential backoff before retrying
		log.Printf("Waiting %v before next attempt...", backoff)
		time.Sleep(backoff)

		// Increase backoff for next attempt (exponential backoff with jitter)
		backoff = time.Duration(float64(backoff) * 1.5)
	}

	return fmt.Errorf("failed to connect to Elasticsearch after %d attempts: %v", maxRetries, lastError)
}

func (s *ESStorage) StorePageData(data PageData) error {
	// Encode the URL to make it safe for use as a document ID
	docID := encodeDocumentID(data.URL)

	// Try to get the existing document by URL (used as ID)
	res, err := s.es.Get(s.index, docID)
	if err != nil {
		return fmt.Errorf("failed to get document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error getting document: %s", res.String())
	}

	// Set timestamps
	now := time.Now()

	// If it exists, merge FullContent if needed and preserve created_at
	if res.StatusCode == 200 {
		var src struct {
			Source PageData `json:"_source"`
		}
		if err := json.NewDecoder(res.Body).Decode(&src); err == nil {
			// Preserve full content if needed
			if src.Source.FullContent != "" && data.FullContent == "" {
				data.FullContent = src.Source.FullContent
			}

			// Preserve created_at timestamp
			data.CreatedAt = src.Source.CreatedAt
			// Update the updated_at timestamp
			data.UpdatedAt = now
		}
	} else {
		// New document, set both timestamps
		data.CreatedAt = now
		data.UpdatedAt = now
	}

	// Index (insert/update) the document
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	res, err = s.es.Index(
		s.index,
		strings.NewReader(string(body)),
		s.es.Index.WithDocumentID(docID),
		s.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (s *ESStorage) GetPage(url string) (*PageData, bool) {
	// Encode the URL to make it safe for use as a document ID
	docID := encodeDocumentID(url)

	res, err := s.es.Get(s.index, docID)
	if err != nil {
		fmt.Printf("Error getting document: %v\n", err)
		return nil, false
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, false
	}
	if res.IsError() {
		fmt.Printf("Error response: %s\n", res.String())
		return nil, false
	}

	var src struct {
		Source PageData `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&src); err != nil {
		fmt.Printf("Error decoding document: %v\n", err)
		return nil, false
	}

	return &src.Source, true
}

func (s *ESStorage) SavePage(page *PageData) error {
	if page == nil {
		return fmt.Errorf("cannot save nil page")
	}
	return s.StorePageData(*page)
}

// NewESStorage creates a new ElasticSearch storage instance
func NewESStorage(client *elasticsearch.Client, indexName string) *ESStorage {
	return &ESStorage{
		es:    client,
		index: indexName,
	}
}

// encodeDocumentID encodes a URL to make it safe for use as a document ID
func encodeDocumentID(url string) string {
	// Use base64 encoding to make the URL safe for use as a document ID
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))
}

// InitializeIndex creates the index if it doesn't exist
func (s *ESStorage) InitializeIndex() error {
	// Check if the index exists
	res, err := s.es.Indices.Exists([]string{s.index})
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %v", err)
	}
	defer res.Body.Close()

	// If the index doesn't exist (404), create it
	if res.StatusCode == 404 {
		// Define index mapping for better search capabilities
		mapping := `{
			"mappings": {
				"properties": {
					"url": { "type": "keyword" },
					"title": { "type": "text", "analyzer": "standard" },
					"snippet": { "type": "text", "analyzer": "standard" },
					"full_content": { "type": "text", "analyzer": "standard" },
					"created_at": { "type": "date" },
					"updated_at": { "type": "date" }
				}
			}
		}`

		// Create the index with the mapping
		res, err := s.es.Indices.Create(
			s.index,
			s.es.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("error creating index: %s", res.String())
		}

		log.Printf("Created Elasticsearch index: %s", s.index)
	} else {
		log.Printf("Elasticsearch index already exists: %s", s.index)
	}

	return nil
}
