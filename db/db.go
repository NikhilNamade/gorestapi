package db

import (
	"database/sql"
	"fmt"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "./data/api.db")
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to open database: %v", err))
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		panic(fmt.Sprintf("❌ Database connection failed: %v", err))
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to enable foreign keys: %v", err))
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	fmt.Println("✅ Database connected successfully.")
	createTables()
}

func createTables() {

	createUserTables := `
		CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(createUserTables)
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to create the user table: %v", err))
	}

	fmt.Println("✅ User table created (or already exists).")
	
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS event (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		date_time DATETIME NOT NULL,
		userId INTEGER,
		profileImage TEXT,
		category TEXT,
		fees INTEGER,
		FOREIGN KEY(userId) REFERENCES user(id) ON DELETE CASCADE
	);
	`

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to create the event table: %v", err))
	}

	createEventsRegistrationTable := `
	CREATE TABLE IF NOT EXISTS register(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	eventId INTEGER,
	userId INTEGER,
	FOREIGN KEY(userId) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY(eventId) REFERENCES event(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createEventsRegistrationTable)
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to create the event registration table: %v", err))
	}

	createStoryTables := `
		CREATE TABLE IF NOT EXISTS Story (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file TEXT NOT NULL,
			userID INTEGER,
			viewStory DATETIME DEFAULT NULL,
			FOREIGN KEY(userID) REFERENCES user(id) ON DELETE CASCADE
		)
	`
	_,err = DB.Exec(createStoryTables);
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to create the story table: %v", err))
	}

	createFollowingTables := `
		CREATE TABLE IF NOT EXISTS Connection (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			followBy INTEGER NOT NULL,
			followTo INTEGER NOT NULL,
			FOREIGN KEY(followBy) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY(followTo) REFERENCES user(id) ON DELETE CASCADE
		)
	`
	_,err = DB.Exec(createFollowingTables);
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to create the connection table: %v", err))
	}
	fmt.Println("✅ Event table created (or already exists).")
}
