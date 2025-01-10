# Makefile

# Variables
BINARY_NAME = ./bin/api
SOURCE = cmd/api/*.go

# Default target
all: build

# Build the project
build:
	@echo "Building the application..."
	@go build -o $(BINARY_NAME) $(SOURCE)

# Run the application
run: build
	@echo "Running the application..."
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(BINARY_NAME)

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Display help
help:
	@echo "Available targets:"
	@echo "  build   - Build the application binary"
	@echo "  run     - Build and run the application"
	@echo "  clean   - Remove build artifacts"
	@echo "  fmt     - Format Go code"
	@echo "  test    - Run all tests"
	@echo "  help    - Display this help message"

