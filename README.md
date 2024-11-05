# Company Service

A Go-based microservice for managing companies, supporting CRUD operations with authentication and optional event streaming.

## Project Overview

This microservice allows clients to create, update, delete, and retrieve details of company records. It is built with best practices for modularity, performance, and production readiness.

### Key Features
- **CRUD Operations**: Supports create, read, update, and delete actions for company records.
- **Authentication**: JWT-based authentication to secure endpoints.
- **Event Streaming**: Kafka-based event handling on data mutations (create, update, delete) (optional).
- **Dockerized**: Easy setup for development and deployment with Docker.
- **Flexible Configurations**: Environment-based configurations for seamless deployments.

## Project Structure

```plaintext
company-service/
├── cmd/                    # Entrypoint for main app
│   └── main.go
├── internal/               # Core application code and business logic
│   ├── company/            # Company CRUD logic and repository
│   ├── auth/               # Authentication logic
│   ├── events/             # Kafka event producer
│   └── db/                 # Database access and repository patterns
├── api/                    # API handlers and routes
│   └── v1/                 # API versioning folder
├── configs/                # Configuration files
│   ├── config.yaml         # Default YAML configuration
│   └── config.go           # Configuration loader
├── scripts/                # Utility scripts
│   └── setup_db.sh         # Script for setting up the database (if needed)
├── tests/                  # Integration and mock tests
│   ├── integration/        # Integration tests for database, API, etc.
│   └── mock/               # Mock objects and utilities for unit tests
├── docker/                 # Docker-related files
│   ├── Dockerfile          # Dockerfile for building the app image
│   ├── docker-compose.yml  # Docker Compose for app and dependencies
│   └── kafka/              # Kafka-specific Docker setup
└── .github/                # GitHub configurations and issue templates
