package company

import (
	config "company-service/configs"
	"company-service/internal/auth"
	"company-service/internal/db"
	"company-service/internal/kafka"
	"log"
	"os"
	"testing"
)

func setupIntegrationTest(t *testing.T) (*CompanyServiceImpl, func()) {
	// Ensure test environment variables are set
	os.Setenv("DATABASE_URL", "postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	os.Setenv("KAFKA_TOPIC_COMPANY_EVENTS", "test_company_events")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load test config: %v", err)
	}

	// Connect to the test database
	database, err := db.Connect() // No argument passed
	if err != nil {
		t.Fatalf("Could not connect to test database: %v", err)
	}

	// Initialize Kafka producer and AuthService
	kafkaProducer := kafka.NewKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopicCompanyEvents)
	authService := auth.NewAuthService(cfg.JWTSecret)

	// Create service
	service := NewCompanyServiceImpl(authService, database, kafkaProducer)

	// Cleanup function
	cleanup := func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing test database: %v", err)
		}
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}

	return service, cleanup
}
