# Makefile for LLM Inference Service

# Variables
BINARY_NAME=llm-inference
GO_VERSION=1.25
PORT=5080
REGISTRY=registry.whatisgoing.com
LOCAL_REGISTRY=localhost:5000
IMAGE_NAME=$(LOCAL_REGISTRY)/$(BINARY_NAME):latest

# Build commands
.PHONY: build
build:
	docker build -t $(BINARY_NAME) .

.PHONY: up
up:
	docker compose up

dev: up bash

.PHONY: down
down:
	docker compose down

.PHONY: bash
bash:
	docker compose exec app bash
.PHONY: clean-images
clean-images:
	docker rmi $(BINARY_NAME) || true

.PHONY: release-local
release-local:
	docker build -t $(BINARY_NAME) .
	docker tag $(BINARY_NAME) $(IMAGE_NAME)
	docker push $(IMAGE_NAME)
	docker push $(REGISTRY)/$(BINARY_NAME):latest
	

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  up            - Run the application"
	@echo "  down          - Stop the application"
	@echo "  clean-images  - Remove the built Docker image"
	@echo "  bash          - Access the running container's bash shell"
	@echo "  help          - Show this help message"
