package main

import (
	"company-service/configs"
	"company-service/internal/auth"
	"company-service/internal/company"
	"company-service/internal/db"
	"company-service/internal/kafka"
	"company-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	authService := auth.NewAuthService(cfg.JWTSecret)

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	kafkaProducer := kafka.NewKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopicCompanyEvents)
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()

	companyService := company.NewCompanyServiceImpl(authService, database, kafkaProducer)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authService.JWTInterceptor),
	)

	proto.RegisterCompanyServiceServer(server, companyService)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server is running on port 8080")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
