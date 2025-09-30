# Local Multilingual NER Service (Go Version)

A high-performance Go service for multilingual Named Entity Recognition (NER) and sentiment analysis using Ollama.

## Features

- **Named Entity Recognition**: Extracts persons, locations, organizations, and events from text
- **Sentiment Analysis**: Analyzes sentiment with confidence scores
- **Multilingual Support**: Arabic, English, and German
- **Fast Performance**: Built with Go and Gin framework
- **Docker Support**: Easy containerization and deployment

## Prerequisites

- Go 1.25 or higher
- Ollama running locally (default: http://localhost:11434)
- Required model: `deepseek-r1:1.5b` (or configure via environment)

## Quick Start

### 1. Install Dependencies
```bash
go mod download
```

### 2. Build and Run
```bash
# Build the binary
make build

# Run directly
make run

# Or run with Go
go run main.go
```

### 3. Using Docker
```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run
```

## API Endpoints
```bash

### NER and Sentiment Analysis
```bash
curl http://localhost:8090/api/v1/text1=some_value&text2=some_value_2&entity_type=persons \
  -H "Content-Type: application/json" \
```

Response:
```
{
  similarity_score: 0.95,
  should_be_merged: true,
  thinking: "some text explaining the llm thinking strategy"
}
```

### Text Similarity

### Health Check
```bash
curl http://localhost:8090/health
```

## Configuration

Set environment variables in `.env` file:
```env
OLLAMA_URL=http://localhost:11434
MODEL_NAME=deepseek-r1:1.5b
PORT=5080
GIN_MODE=release
```

## Build Commands

```bash
# Build for current platform
make build

# Build for Linux
make build-linux

# Build for Windows
make build-windows

# Build for macOS
make build-darwin

# Clean build artifacts
make clean
```

## Network Access

To make the service accessible from other devices on the network, it automatically binds to `0.0.0.0:5080`.

## Development

For development with hot reload:
```bash
go install github.com/codegangsta/gin@latest
make dev
```

## Performance Improvements over Python

- **Faster startup time**: Go compiles to native binary
- **Lower memory usage**: More efficient memory management
- **Better concurrency**: Go's goroutines handle concurrent requests efficiently
- **Single binary deployment**: No dependency management issues
- **Cross-platform builds**: Easy compilation for different platforms