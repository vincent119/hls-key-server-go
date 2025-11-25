.PHONY: all build run test lint fmt tidy clean swagger

# Build variables
BINARY_NAME=hls-key-server
BUILD_DIR=./bin
MAIN_PATH=./cmd/server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

all: fmt lint test build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run:
	@echo "Running application..."
	$(GORUN) $(MAIN_PATH)/main.go

test:
	@echo "Running tests..."
	$(GOTEST) -v -race -count=1 ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

bench:
	@echo "Running benchmarks..."
	$(GOTEST) -run=NONE -bench=. -benchmem ./...

lint:
	@echo "Running linters..."
	$(GOLINT) run ./...

fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	goimports -w .

tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

swagger:
	@echo "Generating Swagger docs..."
	swag init -g $(MAIN_PATH)/main.go -o docs

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 9090:9090 $(BINARY_NAME):latest

help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  bench          - Run benchmarks"
	@echo "  lint           - Run linters"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
