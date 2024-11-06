package db

import (
	config "company-service/configs"
	"database/sql"
	"fmt"
	"log"
	_ "os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
)

func Connect() (*sql.DB, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)

	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping the database: %v", err)
	}

	log.Println("Connected to the database successfully")
	return db, nil
}
