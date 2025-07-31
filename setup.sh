#!/bin/bash

echo "Installing swag for API documentation generation..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "Generating Swagger documentation..."
swag init -g cmd/main.go

echo "Installing dependencies..."
go mod tidy

echo "Building application..."
go build -o bin/crypto-watcher cmd/main.go

echo "Setup complete! You can now run the application with:"
echo "./bin/crypto-watcher"
