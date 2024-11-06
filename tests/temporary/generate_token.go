package main

import (
	"company-service/configs"
	"company-service/internal/auth"
	"log"
)

func main() {
	// Load configuration to get the JWT secret
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	authService := auth.NewAuthService(cfg.JWTSecret)

	token, err := authService.GenerateToken(1)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}

	log.Println("Generated JWT Token:", token)
}
