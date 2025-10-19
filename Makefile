# Makefile for LLM Inference Service

# Variables
BINARY_NAME=llm-inference
GO_VERSION=1.25
PORT=5080
REGISTRY=localhost:5000
IMAGE_NAME=$(REGISTRY)/$(BINARY_NAME):latest

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
	
.PHONY: build
build:
	docker build -t $(BINARY_NAME) .


.PHONY: tag
tag:
	docker tag $(BINARY_NAME) $(IMAGE_NAME)

.PHONY: push
push:
	docker push $(IMAGE_NAME)
.PHONY: release
release:
	make build
	make tag
	make push

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
