package main

import (
	"company-service/configs"
	"company-service/internal/auth"
	"company-service/internal/company"
	"company-service/internal/db"    // Database connection
	"company-service/internal/kafka" // Kafka producer
	"company-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize AuthService with JWT secret
	authService := auth.NewAuthService(cfg.JWTSecret)

	// Connect to the PostgreSQL database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize Kafka producer
	kafkaProducer := kafka.NewKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopicCompanyEvents)
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()

	// Initialize the CompanyService implementation with AuthService, DB, and Kafka producer
	companyService := company.NewCompanyServiceImpl(authService, database, kafkaProducer)

	// Create a new gRPC server with JWT middleware
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authService.JWTInterceptor),
	)

	// Register the CompanyService implementation
	proto.RegisterCompanyServiceServer(server, companyService)

	// Enable gRPC reflection for easier debugging with tools like grpcurl
	reflection.Register(server)

	// Start listening on port 8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server is running on port 8080")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
