.PHONY: build run test clean

# Binary name
BINARY_NAME=dcr-mcp-server

# Build the application
build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with gotestum
test-verbose:
	gotestum --format-hide-empty-pkg --format testdox --format-icons hivis

# Format code
fmt:
	gofumpt -w .

# Clean build files
clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy

# Show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build the application"
	@echo "  run           Run the application"
	@echo "  test          Run tests"
	@echo "  test-verbose  Run tests with verbose output using gotestum"
	@echo "  fmt           Format code using gofumpt"
	@echo "  clean         Clean build files"
	@echo "  deps          Install dependencies"
	@echo "  help          Show this help message"