package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init opens the SQLite database and runs migrations + seeding.
func Init() error {
	dbPath := os.Getenv("RWANDAPI_DB_PATH")
	if dbPath == "" {
		dbPath = "rwandapi.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}

	log.Println("Running database migrations...")
	if err := RunMigrations(DB); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	log.Println("Seeding database...")
	if err := Seed(DB); err != nil {
		return fmt.Errorf("seeding database: %w", err)
	}

	log.Println("Database ready.")
	return nil
}

// Close closes the database connection.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
