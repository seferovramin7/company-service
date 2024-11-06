# Start from the official Go image for building the application
FROM golang:1.23.2 as builder

# Set environment variables for Go modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o company-service ./cmd/main.go

# Start a new stage for the final container
FROM alpine:latest

# Set the working directory inside the final container
WORKDIR /root/

# Copy the built binary from the previous stage
COPY --from=builder /app/company-service .

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["./company-service"]
