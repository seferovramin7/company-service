

# **Company Service Documentation**

## **1. Overview**

This project implements a gRPC-based `company-service` microservice. It provides CRUD operations for company entities and integrates an authentication mechanism with JWT tokens. The service also publishes company-related events to a Kafka topic.

---

## **2. Features**

Non Functional : 
- **gRPC-based microservice** for managing company entities.
- **JWT authentication** to secure gRPC endpoints.
- **Kafka integration** for event-driven architecture.
- **PostgreSQL database** for persistent storage.
- **GitHub Actions CI/CD pipeline** for automated testing and deployment.\

Functional :
- **CRUD Operations**: Supports create, read, update, and delete actions for company records.
- **Authentication**: JWT-based authentication to secure endpoints.
- **Event Streaming**: Kafka-based event handling on data mutations (create, update, delete) (optional).
- **Dockerized**: Easy setup for development and deployment with Docker.
- **Flexible Configurations**: Environment-based configurations for seamless deployments.

---

## **3. Requirements**

- Docker and Docker Compose
- Go 1.20+
- PostgreSQL 14+
- Kafka 7.5+
- grpcurl for testing gRPC endpoints

---

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
```


## **4. Setup Instructions**

### **4.1 Clone the Repository**
```bash
git clone https://github.com/seferovramin7/company-service.git
cd company-service

Start all services using Docker Compose:
docker-compose up --build
```

---

## **5. Testing the Service**

### **5.1 Authentication**

Generate a JWT token:
```bash
grpcurl -plaintext -d '{"user_id": 1}' localhost:8080 company.CompanyService/Login
```

### **5.2 CRUD Operations**

- **Create a Company**:
  ```bash
  grpcurl -plaintext \
    -H "Authorization: Bearer <TOKEN>" \
    -d '{"company": {"name": "Test Co", "description": "A sample company", "employees": 50, "registered": true, "type": "Corporation"}}' \
    localhost:8080 company.CompanyService/CreateCompany
  ```

- **Get a Company**:
  ```bash
  grpcurl -plaintext \
    -H "Authorization: Bearer <TOKEN>" \
    -d '{"id": 1}' \
    localhost:8080 company.CompanyService/GetCompany
  ```

- **Update a Company**:
  ```bash
  grpcurl -plaintext \
    -H "Authorization: Bearer <TOKEN>" \
    -d '{"company": {"id": 1, "name": "Updated Co", "description": "Updated description", "employees": 100, "registered": false, "type": "LLC"}}' \
    localhost:8080 company.CompanyService/UpdateCompany
  ```

- **Delete a Company**:
  ```bash
  grpcurl -plaintext \
    -H "Authorization: Bearer <TOKEN>" \
    -d '{"id": 1}' \
    localhost:8080 company.CompanyService/DeleteCompany
  ```

### **5.3 Verifying Kafka Events**

Look for logs like:
```bash
Successfully published message to Kafka - Key: 1, Message: { "event_type": "CREATE", "company": { ... } }
```

---

## **6. CI/CD Pipeline**

The project uses **GitHub Actions** for CI/CD, including:
- **Automated Tests** for each commit and pull request.
- **Secrets Management** to securely handle sensitive information like `JWT_SECRET` and `DATABASE_URL`.

### **GitHub Actions Workflow**
```yaml
name: CI/CD Pipeline

on:
  push:
    branches:
      - master
  pull_request:
    branches:
      - master
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to deploy (e.g., staging or production)'
        required: false
        default: 'staging'

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up environment variables
        run: |
          echo "JWT_SECRET=${{ secrets.JWT_SECRET }}" >> .env
          echo "DATABASE_URL=${{ secrets.DATABASE_URL }}" >> .env
          echo "KAFKA_BROKER=${{ secrets.KAFKA_BROKER }}" >> .env
          echo "KAFKA_TOPIC_COMPANY_EVENTS=${{ secrets.KAFKA_TOPIC_COMPANY_EVENTS }}" >> .env

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: 1.20

      - name: Install dependencies
        run: go mod tidy

      - name: Run Unit Tests
        run: go test -v ./...

      - name: Build Docker Image
        run: docker build -t company-service .
```

---

## **7. Future Enhancements**

### **7.1 API Gateway**
- A separate **API Gateway** for better authorization handling and routing.

### **7.2 Scaling**
- Load balancing and service discovery for better scalability.

### **7.3 Monitoring**
- Integrate **Prometheus** and **Grafana** for enhanced monitoring and visualization.

---

## **8. Notes**

### **Managing Secrets**
For demonstration purposes, `.env` is included in this repository. In production, always use a secure solution like:
- **GitHub Actions Secrets**
- **AWS Secrets Manager**
- **Vault**

---
```