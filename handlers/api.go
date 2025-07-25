package handlers

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/omed0/go-hello-world/internal/database"
)

var (
	dbConn  *sql.DB
	once    sync.Once
	initErr error
)

type ApiConfig struct {
	Queries *database.Queries
}

// InitDB initializes the database connection and prepares queries.
// It should be called once at application startup.
func (api *ApiConfig) InitDB() error {
	once.Do(func() {
		// Load environment variables only once
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: failed to load .env file: %v", err)
		}

		dbURL := os.Getenv("DB_URL")
		if dbURL == "" {
			initErr = ErrMissingDBURL
			log.Fatal("DB_URL environment variable is required")
			return
		}

		var err error
		dbConn, err = sql.Open("postgres", dbURL)
		if err != nil {
			initErr = err
			log.Fatalf("Cannot open DB: %v", err)
			return
		}

		// Ping DB to verify connection
		if err := dbConn.Ping(); err != nil {
			initErr = err
			log.Fatalf("Cannot reach DB: %v", err)
			return
		}

		api.Queries = database.New(dbConn)
		log.Println("Database connected")
	})

	return initErr
}

// GetDB returns the underlying *sql.DB connection.
// It calls InitDB automatically if not initialized yet.
// If initialization fails, it returns the error.
func (api *ApiConfig) GetDB() (*sql.DB, error) {
	if dbConn == nil {
		err := api.InitDB()
		if err != nil {
			return nil, err
		}
	}
	return dbConn, nil
}

// ErrMissingDBURL is returned when DB_URL is not set.
var ErrMissingDBURL = errors.New("DB_URL environment variable is missing")
