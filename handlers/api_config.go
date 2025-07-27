package handlers

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sync"

	//"time"

	"github.com/joho/godotenv"
	"github.com/omed0/go-hello-world/internal/database"
)

var (
	dbConn  *sql.DB
	once    sync.Once
	initErr error
)

// ErrMissingDBURL is returned when DB_URL is not set.
var ErrMissingDBURL = errors.New("DB_URL environment variable is missing")

type ApiConfig struct {
	Queries *database.Queries
}

// NewApiConfig creates a new ApiConfig instance with a database connection
func NewApiConfig() (*ApiConfig, error) {
	db, err := InstanceDB()
	if err != nil {
		return nil, err
	}

	return &ApiConfig{
		Queries: database.New(db),
	}, nil
}

// InstanceDB initializes the database connection and returns a DB instance.
// It uses a singleton pattern to ensure the database connection is established only once.
func InstanceDB() (*sql.DB, error) {
	once.Do(func() {
		initErr = initializeDB()
	})
	return dbConn, initErr
}

// initializeDB performs the actual database initialization
func initializeDB() error {
	// Load environment variables only once
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return ErrMissingDBURL
	}

	var err error
	dbConn, err = sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	// Configure connection pool
	//dbConn.SetMaxOpenConns(25)
	//dbConn.SetMaxIdleConns(25)
	//dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Ping DB to verify connection
	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		dbConn = nil
		return err
	}

	log.Println("Database connected")
	return nil
}

// CloseDB closes the database connection safely
func CloseDB() {
	if dbConn != nil {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
		dbConn = nil
		// Reset the once to allow reconnection if needed
		once = sync.Once{}
	}
}

// IsDBConnected checks if the database connection is active
func IsDBConnected() bool {
	if dbConn == nil {
		return false
	}
	return dbConn.Ping() == nil
}
