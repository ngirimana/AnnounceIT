package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB // Declare a global variable to hold the database connection

func InitDB() {
	var err error
	// Initialize the DB connection and assign it to the global variable directly
	DB, err = sql.Open("sqlite3", "api.db")
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Set database connection pooling parameters
	DB.SetMaxOpenConns(1000)
	DB.SetMaxIdleConns(950)
	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		last_name TEXT NOT NULL,
		first_name TEXT NOT NULL,
		phone_number TEXT NOT NULL UNIQUE,
		address TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table: " + err.Error())
	}

	createAnnouncementsTable := `
	CREATE TABLE IF NOT EXISTS announcements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_id INTEGER NOT NULL,
		status  string NOT NULL DEFAULT 'pending',
		text TEXT NOT NULL ,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		create_date DATETIME NOT NULL,
		FOREIGN KEY(owner_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createAnnouncementsTable)
	if err != nil {
		panic("Could not create announcements table: " + err.Error())
	}

}

// TruncateUsersTable removes all records from the users table
func TruncateUsersTable() {
	_, err := DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Could not truncate users table: %v", err)
	}
}

func TruncateAnnouncementsTable() {
	_, err := DB.Exec("DELETE FROM announcements")
	if err != nil {
		log.Fatalf("Could not truncate announcements table: %v", err)
	}
}
