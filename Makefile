# Makefile for LLM Inference Service

# Variables
BINARY_NAME=llm-inference
GO_VERSION=1.25
PORT=5080

# Build commands
.PHONY: build
build:
	docker build -t $(BINARY_NAME) .

.PHONY: run
run:
	go run main.go

.PHONY: dev
dev:
	gin --port $(PORT) --appPort $(PORT) --bin ./$(BINARY_NAME) run main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux main.go

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe main.go

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin main.go

# Dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME).exe
	rm -f $(BINARY_NAME)-darwin

# Testing
.PHONY: test
test:
	go test -v ./...

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p $(PORT):$(PORT) --env-file .env $(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-darwin  - Build for macOS"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"