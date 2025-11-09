# Tech newsletter Generator

A Go-based feed collector that monitors various tech blogs and RSS feeds to collect the latest articles.

## Sample

<img width="2186" height="1266" alt="image" src="https://github.com/user-attachments/assets/153df87e-2de5-48cc-b7c2-689bd2dc456b" />

## Project Structure

```
tech-newsletter-generator/
├── cmd/
│   └── collector/          # Feed collector executable
├── internal/
│   ├── config/            # Configuration file management
│   ├── feed/             # Feed fetching and processing
│   └── storage/          # Data storage operations
├── pkg/
│   └── models/           # Data models and structures
├── config/               # Configuration JSON files
├── tmp/data/            # Temporary data storage
└── .github/workflows/   # GitHub Actions workflows
```

## Features

- **Multi-source Feed Support**: Supports RSS, Atom, and various platform-specific feeds (Zenn, Qiita, Hatena, etc.)
- **Intelligent Deduplication**: Avoids collecting duplicate articles
- **Automatic Config Updates**: Updates configuration files with latest article links
- **Comprehensive Testing**: High test coverage with unit tests
- **CI/CD Integration**: GitHub Actions for automated testing and feed collection

## Configuration

Feed configurations are stored in JSON files under the `config/` directory. Each file can contain multiple feed configurations:

```json
{
  "data": [
    {
      "name": "Firebase Blog",
      "type": "categoryIsUrl",
      "feedUrl": "https://firebase.blog/rss.xml",
      "latestLink": "https://firebase.blog/posts/2025/10/fpnv-preview-launch"
    }
  ]
}
```
### GitHub Actions

You have to register GitHub Actions Secret for sending Gmail.

- MAIL_USERNAME
- MAIL_PASSWORD
- MAIL_TO
- MAIL_FROM

`MAIL_PASSWORD` is Gmail App password.

### Supported Feed Types

- `categoryIsUrl`: Direct RSS feed URL
- `categoryIsAtomUrl`: Direct Atom feed URL
- `zenn`: Zenn user feed (username as feedUrl)
- `qiita`: Qiita user feed (username as feedUrl)
- `note`: Note user feed (username as feedUrl)
- `hatena`: Hatena blog feed (blog URL as feedUrl)
- `scrapbox`: Scrapbox project feed (project name as feedUrl)
- `connpass`: Connpass group feed (group name as feedUrl)

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/tech-newsletter-generator.git
   cd tech-newsletter-generator
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Collector

```bash
# Run the feed collector
go run cmd/collector/main.go

# Or use make
make run-collector
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race
```

### Test Coverage

Current test coverage:
- **internal/config**: 87.5%
- **internal/feed**: 93.3%
- **internal/storage**: 84.6%

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Building

```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all
```

## GitHub Actions

### Test Workflow

- Runs on push/PR to master and develop branches
- Tests against Go 1.21, 1.22, and 1.23
- Generates coverage reports
- Uploads coverage to Codecov
- Runs linting and security checks

### Feed Collector Workflow

- Runs every hour automatically
- Can be triggered manually
- Collects new feeds and updates configuration files
- Commits changes back to the repository

## Project Design

### Feed Processing Flow

1. **Load Configurations**: Read all JSON files from `config/` directory
2. **Load Existing Items**: Read `tmp/data/latest-items.json` (create if not exists)
3. **Process Feeds**: For each feed configuration:
   - Fetch latest article from RSS/Atom feed
   - Compare with stored `latestLink` in configuration
   - Check for duplicates in existing items
   - If new article found: update config and add to items
4. **Save Results**: Update configuration files and save new items

### Key Features

- **Date-based Article Detection**: Finds the latest article by publication date, not just the first item
- **Configuration Auto-update**: Prevents duplicate collection by updating `latestLink` in config files
- **Error Resilience**: Continues processing other feeds even if some fail
- **Comprehensive Logging**: Detailed logs for debugging and monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License.
