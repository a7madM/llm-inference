#!/bin/bash

echo "🚀 Setting up Go LLM Inference Service"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Initialize module if needed
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init llm-inference
fi

# Download dependencies
echo "📥 Downloading dependencies..."
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/joho/godotenv@v1.5.1

# Tidy up module
go mod tidy

echo "🔧 Building the service..."
go build -o llm-inference main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎯 To run the service:"
    echo "   ./llm-inference"
    echo ""
    echo "📚 Or use make commands:"
    echo "   make run      # Run directly with go"
    echo "   make build    # Build binary"
    echo "   make help     # Show all commands"
    echo ""
    echo "🌐 Service will be available at: http://localhost:5080"
    echo "📖 API Documentation: http://localhost:5080/"
    echo ""
    echo "🔧 Configuration via .env file:"
    echo "   OLLAMA_URL=http://localhost:11434"
    echo "   MODEL_NAME=deepseek-r1:1.5b"
    echo "   PORT=5080"
else
    echo "❌ Build failed. Please check the error messages above."
    exit 1
fi