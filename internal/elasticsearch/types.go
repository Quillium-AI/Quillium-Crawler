package elasticsearch

import (
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ESStorage struct {
    es      *elasticsearch.Client
    index   string
}

type PageData struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Snippet     string    `json:"snippet"`
	FullContent string    `json:"full_content,omitempty"` // Full HTML content when enabled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
