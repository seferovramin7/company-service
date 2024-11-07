
APP_NAME := company-service
CMD_DIR := ./cmd/main.go
BIN_DIR := ./bin
DOCKER_IMAGE := company-service-image
DOCKER_COMPOSE := docker/docker-compose.yml
KAFKA_COMPOSE := docker/kafka/docker-compose.kafka.yml

.PHONY: all
all: build run

.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

.PHONY: run
run: build
	@echo "Running the application..."
	$(BIN_DIR)/$(APP_NAME)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: lint
lint:
	@echo "Running linting..."
	golangci-lint run

.PHONY: test
test:
	@echo "Running tests with coverage report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	go test -v ./tests/integration

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

.PHONY: docker-up
docker-up:
	@echo "Starting application and dependencies with Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping and removing containers..."
	docker-compose -f $(DOCKER_COMPOSE) down

.PHONY: kafka-up
kafka-up:
	@echo "Starting Kafka services..."
	docker-compose -f $(KAFKA_COMPOSE) up -d

.PHONY: kafka-down
kafka-down:
	@echo "Stopping Kafka services..."
	docker-compose -f $(KAFKA_COMPOSE) down

.PHONY: clean
clean:
	@echo "Cleaning up generated files..."
	rm -rf $(BIN_DIR) coverage.out coverage.html

.PHONY: docs
docs:
	@echo "Generating documentation..."
	go doc > docs/go-docs.md
