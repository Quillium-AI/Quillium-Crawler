FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quillium-crawler .

# Use a smaller image for the final container
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/quillium-crawler .

# Set environment variables with defaults
ENV CRAWLER_ALLOWED_DOMAINS=quilliumtest.com,quilliumexample.com \
    CRAWLER_MAX_DEPTH=3 \
    CRAWLER_THREAD_COUNT=2 \
    CRAWLER_MAX_QUEUE_SIZE=10000 \
    CRAWLER_PARALLELISM=2 \
    CRAWLER_DELAY=1s \
    CRAWLER_RANDOM_DELAY=1s \
    CRAWLER_START_URLS=https://quilliumtest.com,https://quilliumexample.com

# Run the application
CMD ["/app/quillium-crawler"]
