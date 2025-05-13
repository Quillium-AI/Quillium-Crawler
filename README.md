# Quillium-Crawler
The crawler component for Quillium

## Overview
Quillium-Crawler is a web crawler built with Go and the Colly framework. It crawls specified domains and extracts information from web pages.

## Configuration
All configuration is done through environment variables:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| `CRAWLER_ALLOWED_DOMAINS` | Comma-separated list of domains to crawl | quilliumtest.com,quilliumexample.com |
| `CRAWLER_MAX_DEPTH` | Maximum depth to crawl | 3 |
| `CRAWLER_THREAD_COUNT` | Number of threads for the queue | 2 |
| `CRAWLER_MAX_QUEUE_SIZE` | Maximum size of the queue | 10000 |
| `CRAWLER_PARALLELISM` | Number of parallel requests | 2 |
| `CRAWLER_DELAY` | Delay between requests | 1s |
| `CRAWLER_RANDOM_DELAY` | Random delay added to requests | 1s |
| `CRAWLER_START_URLS` | Comma-separated list of URLs to start crawling | https://quilliumtest.com,https://quilliumexample.com |

## Docker Setup

### Using Docker Compose

1. Clone the repository
2. Configure environment variables in `docker-compose.yml` if needed
3. Run the crawler:

```bash
docker-compose up -d
```

### Using Docker Directly

```bash
# Build the Docker image
docker build -t quillium-crawler .

# Run the container with custom configuration
docker run -d \
  --name quillium-crawler \
  -e CRAWLER_ALLOWED_DOMAINS=example.com,example.org \
  -e CRAWLER_MAX_DEPTH=5 \
  -e CRAWLER_START_URLS=https://example.com,https://example.org \
  quillium-crawler
```

## Running Locally

```bash
# Set environment variables (example for Linux/macOS)
export CRAWLER_ALLOWED_DOMAINS=example.com,example.org
export CRAWLER_MAX_DEPTH=5
export CRAWLER_START_URLS=https://example.com,https://example.org

# Run the crawler
go run main.go
```
