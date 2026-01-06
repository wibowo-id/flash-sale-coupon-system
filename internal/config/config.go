package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	// Build database URL from individual variables or use DATABASE_URL if provided
	databaseURL := buildDatabaseURL()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}
}

func buildDatabaseURL() string {
	// If DATABASE_URL is explicitly set, use it (for backward compatibility and Docker)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	// Build from individual database variables
	dbDefault := os.Getenv("DB_DEFAULT")
	if dbDefault == "" {
		dbDefault = "postgresql"
	}

	// Only build PostgreSQL connection string if DB_DEFAULT is postgresql
	if dbDefault == "postgresql" {
		host := os.Getenv("DB_PG_HOST")
		if host == "" {
			host = "localhost"
		}

		database := os.Getenv("DB_PG_DATABASE")
		if database == "" {
			database = "coupon_db"
		}

		username := os.Getenv("DB_PG_USERNAME")
		if username == "" {
			username = "postgres"
		}

		password := os.Getenv("DB_PG_PASSWORD")

		port := os.Getenv("DB_PG_PORT")
		if port == "" {
			port = "5432"
		}

		// Build PostgreSQL connection string
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, username, password, database, port)
	}

	// Default fallback
	return "host=localhost user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable"
}
