package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Build connection string
	var dbURL string
	if strings.HasPrefix(dbHost, "/") {
		// Cloud SQL Unix socket connection (host points to /cloudsql/...)
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("host=%s dbname=%s sslmode=%s", dbHost, dbName, dbSSLMode))
		if dbUser != "" {
			builder.WriteString(fmt.Sprintf(" user=%s", dbUser))
		}
		if dbPassword != "" {
			builder.WriteString(fmt.Sprintf(" password=%s", dbPassword))
		}
		if dbPort != "" {
			builder.WriteString(fmt.Sprintf(" port=%s", dbPort))
		}
		dbURL = builder.String()
		fmt.Printf("🔌 Using Cloud SQL socket connection: host=%s\n", dbHost)
	} else if dbPassword != "" {
		// TCP connection with password
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
		fmt.Printf("🔌 Using TCP connection: postgres://%s:***@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	} else {
		// TCP connection without password (e.g., local dev)
		dbURL = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
			dbUser, dbHost, dbPort, dbName, dbSSLMode)
		fmt.Printf("🔌 Using TCP connection without password: postgres://%s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	}

	// Open database connection
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
