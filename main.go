package main

import (
	"database/sql"
	"log"
)

// Represent DB connection
var (
	DB             *sql.DB
	target_db_file string
)

// Initalize db
func InitDB() {
	target_db_file = "./spend.db"
	var err error

	DB, err = sql.Open("sqlite3", target_db_file)
	asser_error("Error initializing DB: %q", err)

	init_schema := `
	CREATE TABLE IF NOT EXISTS daily(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		item TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		tag_id INTEGER,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS monthly(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		item TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		tag_id INTEGER,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS tags(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS spend_ahead(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS forecast(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		daily_mean REAL NOT NULL,
		daily_total REAL NOT NULL,
		daily_forecast REAL NOT NULL,
		monthly_total REAL NOT NULL,
		ahead_total REAL NOT NULL,
		income_total REAL NOT NULL,
		overshoot_total REAL NOT NULL
	);
	`
	_, err = DB.Exec(init_schema)
	asser_error("Error executing init DB: %q", err)
}

func Create(input [5]string) [5]string {
	return input
}

// log error
func asser_error(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func main() {
	InitDB()
}
