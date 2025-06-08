# Contributing to Quillium

## IMPORTANT: GitLab is our primary repository

**PLEASE NOTE:** This GitHub repository is a mirror. All contributions should be made through our GitLab repository:
- Repository URL: https://gitlab.cherkaoui.ch/quillium-ai/quillium-crawler
- Issue Board: https://gitlab.cherkaoui.ch/quillium-ai/quillium-crawler/-/issues

**Any pull requests or issues opened on GitHub will be ignored and closed after a warning.**

Please fork and create merge requests on GitLab instead.

Thank you for your interest in contributing to Quillium! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## CLA (Contributor License Agreement)

Before your contribution can be accepted, you need to sign our Contributor License Agreement (CLA). We use CLA Assistant (cla-assistant.io) to manage our CLA process.

When you create a pull request, CLA Assistant will automatically check if you've signed the CLA. If not, you'll be prompted to do so directly within the pull request by clicking on a link that will take you to the CLA Assistant website where you can sign the agreement with your GitHub account.

## Commit Message Conventions

We follow standard git commit message conventions as outlined in our [Code of Conduct](CODE_OF_CONDUCT.md#commit-message-conventions). Please ensure your commit messages follow this format.

## Documentation

Comprehensive documentation for the project is available at [docs.quillium.dev](https://docs.quillium.dev). Please refer to the documentation for detailed information about:

- [API Reference](https://docs.quillium.dev/backend/api/)
- [Authentication System](https://docs.quillium.dev/backend/authentication/)
- [Database Schema](https://docs.quillium.dev/backend/database/)
- [Testing Guidelines](https://docs.quillium.dev/backend/testing/)

When contributing new features or making significant changes, please update the relevant documentation as well.

## How to Contribute

### Reporting Bugs

Bugs are tracked as GitHub issues. Search the [issues](https://github.com/Quillium-AI/Quillium/issues) to see if your bug has already been reported. If not, create a new issue with a clear description and as much relevant information as possible.

### Suggesting Enhancements

Enhancement suggestions are also tracked as GitHub issues. Please provide clear descriptions of the enhancement and how it would benefit the project.

### Pull Requests

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes using the conventional commit format (`git commit -m 'feat(component): add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Workflow

To contribute to Quillium-Crawler, you'll need to set up your local backend environment. Please follow these steps:

1. Clone the repository
   ```bash
   git clone https://github.com/Quillium-AI/Quillium-Crawler.git
   cd Quillium-Crawler
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Configure your environment
   ```bash
   # Copy the example environment file
   cp .env.example .env
   # Edit the .env file to configure your crawler
   ```

4. Run the crawler
   ```bash
   go run main.go
   ```

3. Make your changes and ensure they pass all tests
   ```bash
   # Run tests
   go test ./...
   ```

### Docker Development

To run the crawler using Docker:

```bash
# Build and run with Docker Compose
docker compose up -d --build

# Or build and run manually
docker build -t quillium-crawler .
docker run --env-file .env quillium-crawler
```

The Docker setup uses the environment variables from your `.env` file, so make sure it's properly configured before building.

## Project Structure

The Quillium-Crawler project is organized as follows:

```
Quillium-Crawler/
├── internal/
│   ├── api/              # API server and routes
│   ├── crawler/          # Crawling logic, config, anti-bot, proxies
│   ├── dedup/            # Deduplication (bloom filter)
│   ├── elasticsearch/    # Elasticsearch integration
│   └── metrics/          # Metrics and monitoring
├── main.go               # Application entry point
├── Dockerfile            # Docker build file
├── .env.example          # Example environment variables
├── README.md             # Project overview
├── CONTRIBUTING.md       # Contribution guidelines
├── go.mod                # Go module definition
├── go.sum                # Go dependency checksums
├── docker-compose.yml    # Docker Compose configuration
```

For more detailed information about the project structure and architecture, please refer to the [documentation](https://docs.quillium.dev).

## Style Guidelines

### Code Style

- **Go**: Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and use `gofmt` to format your code.
- **Go**: Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and use `gofmt` to format your code.

### Testing Requirements

- All new features should include appropriate unit tests and integration tests using Go's testing package.

### Commit Messages

Follow the commit message format as described in the [Code of Conduct](CODE_OF_CONDUCT.md#commit-message-conventions):

```
<type>(<scope>): <description>
```

For example:
- `feat(auth): add login functionality`
- `fix(api): resolve null pointer in user data fetch`

## License

By contributing to Quillium, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
