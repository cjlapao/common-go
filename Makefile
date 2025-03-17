# Makefile for common-go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
BINARY_NAME=common-go
VERSION=$(shell cat VERSION_FILE 2>/dev/null || echo "dev")

# Linting tools
GOLINT=golangci-lint
GOSEC=gosec

# Security tools
GOVULNCHECK=govulncheck

# Build flags
LDFLAGS=-ldflags "-X github.com/cjlapao/common-go/version.Version=$(VERSION)"

.PHONY: all build clean test lint security deps tidy help update-deps

all: deps build test lint security

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linting
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Run vet
vet:
	@echo "Running vet..."
	$(GOVET) ./...

# Run security checks
security:
	@echo "Running security checks..."
	$(GOSEC) ./...
	$(GOVULNCHECK) ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/securego/gosec/v2/cmd/gosec
	$(GOGET) -u golang.org/x/vuln/cmd/govulncheck

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Help
help:
	@echo "Make targets:"
	@echo "  all          - Build, test, lint, and run security checks"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint         - Run linter"
	@echo "  vet          - Run go vet"
	@echo "  security     - Run security checks"
	@echo "  deps         - Install dependencies"
	@echo "  update-deps  - Update dependencies"
	@echo "  tidy         - Tidy go.mod"
	@echo "  fmt          - Format code"
	@echo "  help         - Show this help message" 