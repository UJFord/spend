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

// Initialize db
func InitDB() sql.Result {
	target_db_file = "./spend.db"
	var result sql.Result
	var err error

	DB, err = sql.Open("sql3", target_db_file)
	assert_error("Error opening sql.Open: %q", err)

	init_table := `
	CREATE TABLE IF NOT EXISTS daily(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		item TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS tags(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS forecast(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		daily_mean REAL NOT NULL,
		daily_total REAL NOT NULL,
		monthly_total REAL NOT NULL,
		daily_forecast REAL NOT NULL,
		daily_overshoot REAL NOT NULL,
		spend_total REAL NOT NULL,
		income_total REAL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS monthly(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);`

	result, err = DB.Exec(init_table)
	assert_error("Error executing DB.Exec: %q", err)

	return result
}

func Create(input [5]string) [5]string {
	return input
}

// Assert Error
func assert_error(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func main() {
	InitDB()
}
