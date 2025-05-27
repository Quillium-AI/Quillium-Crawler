# Quillium-Crawler

Quillium-Crawler is a high-performance, extensible web crawler written in Go. It supports advanced anti-bot evasion, flexible domain/URL filtering, metrics, and Elasticsearch integration. For full documentation, see [docs.quillium.dev](https://docs.quillium.dev).


## Quick Start

1. Copy `.env.example` to `.env` and adjust as needed.
2. Run with Docker or locally (see below).

## Key Environment Variables
See `.env.example` for all options. Common variables:

| Variable                        | Example / Default         | Description                   |
|----------------------------------|--------------------------|-------------------------------|
| CRAWLER_START_URLS               | https://example.com      | Comma-separated start URLs    |
| CRAWLER_ALLOWED_DOMAINS          | example.com,example2.com | Domains to crawl              |
| CRAWLER_MAX_DEPTH                | 3                        | Max crawl depth               |
| CRAWLER_PARALLEL_REQUESTS        | 10                       | Parallel requests             |
| CRAWLER_ENABLE_METRICS           | true                     | Enable Prometheus metrics     |
| CRAWLER_INDEX_NAME               | crawled_data             | Elasticsearch index           |
| CRAWLER_ENABLE_USER_AGENT_ROTATION| true                    | Anti-bot: rotate user agents  |
| CRAWLER_ENABLE_HEADER_RANDOMIZATION| true                   | Anti-bot: randomize headers   |
| CRAWLER_ELASTICSEARCH_ADDRESSES  | http://localhost:9200    | Elasticsearch endpoint        |

Refer to `.env.example` for advanced anti-bot and proxy settings.

### Anti-Bot Features

- User agent rotation
- Header randomization
- Cookie handling
- Sophisticated/randomized delays

Enable/disable via environment variables. See `.env.example` for details.

## Running with Docker

```bash
docker-compose up -d
```

Or build and run manually:

```bash
docker build -t quillium-crawler .
docker run --env-file .env quillium-crawler
```

## Local Development

1. Copy `.env.example` to `.env` and adjust as needed.
2. Install Go (>=1.20).
3. Download dependencies:
   ```bash
   go mod download
   ```
4. Run the crawler:
   ```bash
   go run main.go
   ```

For Elasticsearch, ensure your ES instance is running and credentials are set in your `.env`.
