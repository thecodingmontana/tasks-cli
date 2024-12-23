package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// ConnectDB initializes the database connection and creates the tasks table if it doesn't exist.
func ConnectDB() {
	var err error
	db, err = sql.Open("sqlite3", "./pkg/database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create tasks table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);	
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// GetDB provides the database connection instance.
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}
}
