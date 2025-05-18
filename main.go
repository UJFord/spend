package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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
	assert_error("Error initializing DB", err)

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
	CREATE TABLE IF NOT EXISTS ahead(
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
	assert_error("Error executing init DB", err)
}

// Insert a spend
func CreateDaily(input []string) ([]string, error) {
	item := input[0]
	amount := input[1]
	date := input[2]
	tag := input[3]

	insert_stmt, err := DB.Prepare(`
		INSERT INTO daily(item, amount, date, tag)
		VALUES (?, ?, ?, ?)
	`)
	assert_error("Error preparing insert statement", err)
	defer insert_stmt.Close()

	return input, err
}

// Inserting or getting a tag
func tag_get_or_insert(tag_name string) int64 {
	var tag_id int64
	err := DB.QueryRow("SELECT id FROM tags WHERE name = ?", tag_name).Scan(&tag_id)

	if err == sql.ErrNoRows {
		statement, err := DB.Prepare("INSERT INTO tags(name) VALUES (?)")
		assert_error(fmt.Sprintf("Error 'preparing for insert' new tag(%s) in tags table", tag_name), err)
		defer statement.Close()

		result, err := statement.Exec(tag_name)
		assert_error(fmt.Sprintf("Error 'inserting' new tag(%s) in tags table", tag_name), err)

		last_insert_id, err := result.LastInsertId()
		assert_error(fmt.Sprintf("Error getting LastInsertId after inserting new tag(%s)", tag_name), err)

		return last_insert_id
	}

	assert_error(fmt.Sprintf("Error querying row of '%s'", tag_name), err)
	return tag_id

}

// Log error
func assert_error(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %q", message, err)
	}
}

// Validate input
func validate_input(args []string) {
	action := args[0]

	switch action {
	case "-cd":
		CreateDaily(args[1:])
	default:
		assert_error(fmt.Sprintf("Invalid action '%s'", action), errors.New("Action not recognized"))
	}
}

func main() {
	// InitDB()
	fmt.Println(os.Args[1:])
}
