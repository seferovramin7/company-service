package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	JWTSecret               string
	AppPort                 string
	DatabaseURL             string
	KafkaBroker             string
	KafkaTopicCompanyEvents string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env") // Optional if you have an .env file
	viper.AutomaticEnv()        // Automatically read environment variables

	viper.SetDefault("JWT_SECRET", "mySecretKey")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("KAFKA_BROKER", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC_COMPANY_EVENTS", "company_events")

	err := viper.ReadInConfig() // Optional: Reads from .env if available
	if err != nil {
		log.Printf("Config file not found, using environment variables instead")
	}

	config := &Config{
		JWTSecret:               viper.GetString("JWT_SECRET"),
		AppPort:                 viper.GetString("APP_PORT"),
		DatabaseURL:             viper.GetString("DATABASE_URL"),
		KafkaBroker:             viper.GetString("KAFKA_BROKER"),
		KafkaTopicCompanyEvents: viper.GetString("KAFKA_TOPIC_COMPANY_EVENTS"),
	}

	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	return config, nil
}
