FROM golang:1.23.2 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o company-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/company-service .

EXPOSE 8080

CMD ["./company-service"]
