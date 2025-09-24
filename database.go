package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDatabase() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Build connection string
	var dbURL string
	if dbPassword != "" {
		// TCP connection with password
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
		fmt.Printf("🔌 Using TCP connection: postgres://%s:***@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	} else {
		// Unix socket connection (no password needed)
		dbURL = fmt.Sprintf("postgres://%s@/%s?sslmode=%s&host=/var/run/postgresql",
			dbUser, dbName, dbSSLMode)
		fmt.Printf("🔌 Using Unix socket connection\n")
	} // Open database connection
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database!")
}

func closeDatabase() {
	if db != nil {
		db.Close()
	}
}
