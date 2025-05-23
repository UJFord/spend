package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
func CreateDaily(input []string) (string, int64) {
	item := input[0]
	amount, err := strconv.Atoi(input[1])
	assert_error(fmt.Sprintf("Error converting %s into int", input[1]), err)
	date := validate_date(input[2])
	tag := tag_get_or_insert(input[3])

	insert_stmt, err := DB.Prepare(`
		INSERT INTO daily(item, amount, date, tag_id)
		VALUES (?, ?, ?, ?)
	`)
	assert_error("Error preparing insert statement", err)
	defer insert_stmt.Close()

	exec_insert_stmt, err := insert_stmt.Exec(item, amount, date, tag)
	assert_error("Error executing insert statement", err)

	id_of_inserted, err := exec_insert_stmt.LastInsertId()
	assert_error("Error fetching last insert id", err)

	output := fmt.Sprintf("Daily Spend Created: %s with id %d\n", strings.Join(input, " "), id_of_inserted)

	return output, id_of_inserted
}

// Edit a daily spend
func EditDaily(target_daily_id int64, target_info int, to_replace_with string) (string, string) {

	if target_info < 0 || target_info > 3 {
		log.Fatalf("Error choosing target info: Out of Bounds")
	}

	var target string
	switch target_info {
	case 0:
		target = "item"
	case 1:
		target = "amount"
	case 2:
		target = "date"
	case 3:
		target = "tag_id"
		to_replace_with = strconv.FormatInt(tag_get_or_insert(to_replace_with), 10)
	}

	update_stmt, err := DB.Prepare(fmt.Sprintf("UPDATE daily SET %s = ? WHERE id = ?", target))
	assert_error("Error preparing edit statement", err)
	defer update_stmt.Close()

	exec_update_stmt, err := update_stmt.Exec(to_replace_with, target_daily_id)
	assert_error("Error executing update statement", err)

	inserted_id, err := exec_update_stmt.LastInsertId()
	replaced_value := get_daily_by_id(inserted_id)[target_info]

	return fmt.Sprintf("Edited Daily Spend: %d from %s into %s", target_daily_id, replaced_value, to_replace_with), to_replace_with
}

// Remove a daily spend
func RemoveDaily(target_id int64) string {

	target_daily := get_daily_by_id(target_id)

	delete_stmt, err := DB.Prepare("DELETE FROM daily WHERE id=?")
	assert_error("Error preparing delete statement:", err)

	_, err = delete_stmt.Exec(target_id)
	assert_error("Error executing delete statement:", err)

	return fmt.Sprintf("Removed Daily Spend: %d %s", target_id, strings.Join(target_daily[:], " "))
}

// Get Daily info by id
func get_daily_by_id(target_id int64) [4]string {

	get_daily := DB.QueryRow("SELECT item, amount, date, tag_id FROM daily WHERE id=?", target_id)

	var result [4]string
	err := get_daily.Scan(&result[0], &result[1], &result[2], &result[3])
	assert_error("Error scanning Get Daily Info by ID statement", err)

	result[2] = get_date_from_time_struct(result[2])

	result[3] = get_tag_by_id(result[3])

	return result
}

//

// Get Date from time.Time structure
func get_date_from_time_struct(time_struct string) string {

	t, err := time.Parse("2006-01-02 15:04:05-07:00", time_struct)
	assert_error(fmt.Sprintf("Error Parsing time.Time %s", time_struct), err)

	date := fmt.Sprintf("%d-%d-%d", int(t.Month()), t.Day(), t.Year())

	return date

}

// Get tag name by id
func get_tag_by_id(target_id string) string {

	get_tag_name := DB.QueryRow("SELECT name FROM tags WHERE id=?", target_id)

	var result string
	err := get_tag_name.Scan(&result)
	assert_error("Error scanning Get Tag Name by ID", err)

	return result

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

// Date formatting
func validate_date(unparsed string) time.Time {
	split_unparsed := strings.Split(unparsed, "-")

	month, err := strconv.Atoi(split_unparsed[0])
	assert_error("Error converting month to int", err)

	day, err := strconv.Atoi(split_unparsed[1])
	assert_error("Error converting day to int", err)

	year, err := strconv.Atoi(split_unparsed[2])
	assert_error("Error converting year to int", err)

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if !(t.Year() == year && int(t.Month()) == month && t.Day() == day) {
		log.Fatalf("Invalid date %d-%d-%d", month, day, year)
	}

	return t
}

// Log error
func assert_error(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %q", message, err)
	}
}

// Validate input
func Validate(args []string) (string, int64) {
	action := args[0]

	switch action {
	case "-cd":
		return CreateDaily(args[1:])
	default:
		assert_error(fmt.Sprintf("Invalid action '%s'", action), errors.New("Action not recognized"))
	}

	return "There should be an error", 0
}

func main() {
	// fmt.Println(os.Args[1:])
	InitDB()
	fmt.Println(Validate(os.Args[1:]))
	// fmt.Println(get_daily_by_id(1))
}
