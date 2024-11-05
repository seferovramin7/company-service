# Define application binary name and directories
APP_NAME := company-service
CMD_DIR := ./cmd/main.go
BIN_DIR := ./bin
DOCKER_IMAGE := company-service-image
DOCKER_COMPOSE := docker/docker-compose.yml
KAFKA_COMPOSE := docker/kafka/docker-compose.kafka.yml

# Default command - builds and runs the application
.PHONY: all
all: build run

# Clean up old binaries and build the application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

# Run the application binary
.PHONY: run
run: build
	@echo "Running the application..."
	$(BIN_DIR)/$(APP_NAME)

# Format all Go files in the project
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linting to enforce code quality standards
.PHONY: lint
lint:
	@echo "Running linting..."
	golangci-lint run

# Run unit tests with coverage report
.PHONY: test
test:
	@echo "Running tests with coverage report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests (assuming tests are in the ./tests/integration directory)
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	go test -v ./tests/integration

# Docker: Build the Docker image for the application
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

# Docker: Run the application and external services with Docker Compose
.PHONY: docker-up
docker-up:
	@echo "Starting application and dependencies with Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) up -d

# Docker: Stop and remove containers created by Docker Compose
.PHONY: docker-down
docker-down:
	@echo "Stopping and removing containers..."
	docker-compose -f $(DOCKER_COMPOSE) down

# Docker: Run only Kafka setup for local development
.PHONY: kafka-up
kafka-up:
	@echo "Starting Kafka services..."
	docker-compose -f $(KAFKA_COMPOSE) up -d

# Docker: Stop Kafka services
.PHONY: kafka-down
kafka-down:
	@echo "Stopping Kafka services..."
	docker-compose -f $(KAFKA_COMPOSE) down

# Clean up binaries and other generated files
.PHONY: clean
clean:
	@echo "Cleaning up generated files..."
	rm -rf $(BIN_DIR) coverage.out coverage.html

# Generate documentation (example for go doc)
.PHONY: docs
docs:
	@echo "Generating documentation..."
	go doc > docs/go-docs.md
