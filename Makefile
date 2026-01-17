.PHONY: build test fmt lint clean install

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
all: build

# Build the binary
build:
	go build $(LDFLAGS) -o bin/ynabctl .

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Run linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/

# Install locally
install: build
	cp bin/ynabctl $(GOPATH)/bin/ynabctl

# Run go mod tidy
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Verify dependencies
verify:
	go mod verify

# Build for all platforms (for local testing)
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/ynabctl-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/ynabctl-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/ynabctl-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/ynabctl-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/ynabctl-windows-amd64.exe .

# Run the application
run:
	go run $(LDFLAGS) . $(ARGS)

# Check code (fmt + lint + test)
check: fmt lint test
