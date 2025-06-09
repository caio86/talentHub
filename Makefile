# Makefile for Go project

# Customize these variables for your project
BINARY_NAME := talentHub
GO_PACKAGE := ./cmd/talentHub # Path to main package
# VERSION ?= $(shell git describe --tags || echo "v0.0.0")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOARCH ?= amd64
GOOS ?= $(shell go env GOOS)

# Default target
all: build

# Build the project
build:
	go build -o bin/$(BINARY_NAME) $(GO_PACKAGE)

# Cross-compilation targets
build-linux:
	GOOS=linux GOARCH=$(GOARCH) go build -o bin/$(BINARY_NAME)-linux-$(GOARCH) $(GO_PACKAGE)

build-darwin:
	GOOS=darwin GOARCH=$(GOARCH) go build -o bin/$(BINARY_NAME)-darwin-$(GOARCH) $(GO_PACKAGE)

build-windows:
	GOOS=windows GOARCH=$(GOARCH) go build -o bin/$(BINARY_NAME)-windows-$(GOARCH).exe $(GO_PACKAGE)

# Run the application
run:
	go run $(GO_PACKAGE)

# Clean build artifacts
clean:
	go clean
	rm -rf bin/ cover.out

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

swag:
	swag i -g ./http/server.go

# Show help
help:
	@echo "Available targets:"
	@echo "  all       - Default build target (alias: build)"
	@echo "  build     - Build binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-darwin  - Build for macOS"
	@echo "  build-windows - Build for Windows"
	@echo "  run       - Run application directly"
	@echo "  clean     - Remove build artifacts"
	@echo "  test      - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt       - Format source code"
	@echo "  vet       - Analyze source code"
	@echo "  deps      - Download dependencies"
	@echo "  help      - Show this help"

.PHONY: all build build-linux build-darwin build-windows run clean test test-coverage fmt vet deps help
