# Quillium-Crawler
The crawler component for Quillium

## Overview
Quillium-Crawler is a web crawler built with Go and the Colly framework. It crawls specified domains and extracts information from web pages.

## Configuration
All configuration is done through environment variables:

| Environment Variable          | Description                                                 | Default Value                                        |
| ----------------------------- | ----------------------------------------------------------- | ---------------------------------------------------- |
| `CRAWLER_ALLOWED_DOMAINS`     | Comma-separated list of domains to crawl                    | quilliumtest.com,quilliumexample.com                 |
| `CRAWLER_MAX_DEPTH`           | Maximum depth to crawl                                      | 3                                                    |
| `CRAWLER_THREAD_COUNT`        | Number of threads for the queue                             | 2                                                    |
| `CRAWLER_MAX_QUEUE_SIZE`      | Maximum size of the queue                                   | 10000                                                |
| `CRAWLER_PARALLELISM`         | Number of parallel requests                                 | 2                                                    |
| `CRAWLER_DELAY`               | Delay between requests                                      | 1s                                                   |
| `CRAWLER_RANDOM_DELAY`        | Random delay added to requests                              | 1s                                                   |
| `CRAWLER_START_URLS`          | Comma-separated list of URLs to start crawling              | https://quilliumtest.com,https://quilliumexample.com |
| `CRAWLER_ACCEPT_LANGUAGE`     | Value for Accept-Language HTTP header (e.g. en-US,en;q=0.9) | (none)                                               |
| `CRAWLER_ENABLE_FULL_CONTENT` | Enable full page content scraping                           | false                                                |
| `CRAWLER_ENABLE_METRICS`      | Enable Prometheus metrics endpoint                          | false                                                |

### Anti-Bot Configuration

Quillium-Crawler includes several anti-bot detection measures to help avoid being blocked by websites:

| Environment Variable                  | Description                                                | Default Value |
| ------------------------------------- | ---------------------------------------------------------- | ------------- |
| `CRAWLER_ENABLE_USER_AGENT_ROTATION`  | Enable random user agent rotation                          | true          |
| `CRAWLER_ENABLE_HEADER_RANDOMIZATION` | Enable HTTP header randomization                           | true          |
| `CRAWLER_ENABLE_COOKIE_HANDLING`      | Enable browser-like cookie handling                        | true          |
| `CRAWLER_ENABLE_SOPHISTICATED_DELAYS` | Enable more human-like delays                              | true          |
| `CRAWLER_RANDOM_DELAY_FACTOR`         | Factor for random delay calculation                        | 1.5           |
| `CRAWLER_CUSTOM_USER_AGENTS`          | Additional custom user agents (comma-separated)            | (none)        |
| `CRAWLER_CUSTOM_ACCEPT_LANGUAGES`     | Additional custom accept-language values (comma-separated) | (none)        |

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
export CRAWLER_ACCEPT_LANGUAGE="en-US,en;q=0.9"

# Run the crawler
go run main.go
```
