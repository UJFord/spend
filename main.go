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
	CREATE TABLE IF NOT EXISTS spend(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		item TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT DEFAULT '1970-01-01 00:00:00+00:00',
		tag_id INTEGER,
		is_daily INTEGER DEFAULT '1',
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS tags(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS ahead(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		date TEXT DEFAULT '1970-01-01 00:00:00+00:00'
	);
	CREATE TABLE IF NOT EXISTS income(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT DEFAULT '1970-01-01 00:00:00+00:00'
	);
	`
	_, err = DB.Exec(init_schema)
	assert_error("Error executing init DB", err)
}

// Create
func Create(input []string) (string, int64) {

	item := input[0]

	amount, err := strconv.Atoi(input[1])
	assert_error(fmt.Sprintf("CREATE error converting %s into int", input[2]), err)

	date := ParseDate(input[2])

	tag := GetTagID(input[3])

	var is_daily bool
	switch input[4] {
	case "":
		is_daily = true
	case "daily":
		is_daily = true
	case "monthly":
		is_daily = false
	default:
		log.Fatalf("CREATE don't know what %s means", input[4])
	}

	insert_stmt, err := DB.Prepare(fmt.Sprintf(`
		INSERT INTO spend(item, amount, date, tag_id, is_daily)
		VALUES (?, ?, ?, ?, ?)
	`))
	assert_error("CREATE error preparing insert statement", err)
	defer insert_stmt.Close()

	exec_insert_stmt, err := insert_stmt.Exec(item, amount, date, tag, is_daily)
	assert_error("CREATE error executing insert statement", err)

	id_of_inserted, err := exec_insert_stmt.LastInsertId()
	assert_error("CREATE error fetching last insert id", err)

	output := fmt.Sprintf("CREATE spend created: %s with id %d\n", strings.Join(input, " "), id_of_inserted)

	return output, id_of_inserted
}

// Get info by id
func Read(target_id int64) ([5]string, string) {

	get := DB.QueryRow("SELECT item, amount, date, tag_id, is_daily FROM spend WHERE id=?", target_id)

	var result [5]string
	err := get.Scan(&result[0], &result[1], &result[2], &result[3], &result[4])
	assert_error(fmt.Sprintf("READ error scanning get spend info by id(%d) statement", target_id), err)

	result[2] = UnparseDate(result[2])

	result[3] = GetTagValue(result[3])

	var is_daily string
	switch result[4] {
	case "1":
		is_daily = "daily"
	case "0":
		is_daily = "monthly"
	default:
		log.Fatal()
	}
	result[4] = is_daily

	output := fmt.Sprintf("READ spend info: %d %s", target_id, strings.Join(result[:], " "))

	return result, output
}

// Edit a spend
func Edit(target_id int64, target_info int, to_replace_with string) (string, string) {

	if target_info < 0 || target_info > 3 {
		log.Fatalf("EDIT error choosing target info: Out of Bounds")
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
		to_replace_with = strconv.FormatInt(GetTagID(to_replace_with), 10)
	}

	edit_stmt, err := DB.Prepare(fmt.Sprintf("UPDATE spend SET %s = ? WHERE id = ?", target))
	assert_error("EDIT error preparing edit statement", err)
	defer edit_stmt.Close()

	_, err = edit_stmt.Exec(to_replace_with, target_id)
	assert_error("EDIT error executing edit statement", err)

	daily_value, _ := Read(target_id)
	replaced_value := daily_value[target_info]

	return fmt.Sprintf("EDIT edited spend: id(%d) from (%s) into (%s)",
			target_id, replaced_value, to_replace_with),
		to_replace_with
}

// Remove spend
func Remove(target_id int64) string {

	target, _ := Read(target_id)

	delete_stmt, err := DB.Prepare("DELETE FROM spend WHERE id=?")
	assert_error("REMOVE error preparing delete statement:", err)

	_, err = delete_stmt.Exec(target_id)
	assert_error("REMOVE error executing delete statement:", err)

	return fmt.Sprintf("REMOVE removed spend: %d %s", target_id, strings.Join(target[:], " "))
}

// Create spend ahead
func CreateAhead(amount float64, date string) (string, int64) {

	parsed_date := ParseDate(date)

	insert_stmt, err := DB.Prepare("INSERT INTO ahead(amount, date) VALUES(?, ?)")

	assert_error("CREATE error preparing insert statement", err)
	defer insert_stmt.Close()

	exec_insert_stmt, err := insert_stmt.Exec(amount, parsed_date)
	assert_error("CREATE error executing insert statement", err)

	id_of_inserted, err := exec_insert_stmt.LastInsertId()
	assert_error("CREATE error fetching last insert id", err)

	date = parsed_date.Format("1-2-2006")
	output := fmt.Sprintf("CREATE AHEAD spending amount(%.2f) date(%s) ahead with id(%d)",
		amount,
		date,
		id_of_inserted)

	return output, id_of_inserted

}

// Read spend ahead
func ReadAhead(target_id int64) (float64, string) {

	get := DB.QueryRow("SELECT amount, date FROM ahead WHERE id=?", target_id)

	var result [2]string
	err := get.Scan(&result[0], &result[1])
	assert_error(fmt.Sprintf("READ AHEAD error scanning get item info by id(%d) statement", target_id), err)

	parsed_amount, err := strconv.ParseFloat(result[0], 64)
	assert_error("READ AHEAD error parsing float", err)
	parsed_date := UnparseDate(result[1])

	output := fmt.Sprintf("READ AHEAD id(%d) amount(%.2f) date(%s)", target_id, parsed_amount, parsed_date)

	return parsed_amount, output

}

// Remove spend ahead
func RemoveAhead(target_id int64) string {

	amount, _ := ReadAhead(target_id)

	delete_stmt, err := DB.Prepare("DELETE FROM ahead WHERE id=?")
	assert_error("REMOVE AHEAD error preparing delete statement:", err)

	_, err = delete_stmt.Exec(target_id)
	assert_error("REMOVE AHEAD error executing delete statement:", err)

	return fmt.Sprintf("REMOVE AHEAD spending amount(%.2f) ahead with id(%d)", amount, target_id)
}

// Forecast
func Forecast() {
	// ahead sum
	get_total := func(table string) (float64, error) {

		query := fmt.Sprintf(`
			SELECT SUM(amount) FROM %s
			WHERE strftime('%%Y-%%m', date) = strftime('%%Y-%%m', now)
		`, table)

		var total sql.NullFloat64
		err := DB.QueryRow(query).Scan(&total)
		if err != nil {
			return 0, err
		}

		if total.Valid {
			return total.Float64, nil
		}
		return 0, nil
	}

	get_mean := func() (float64, error) {
		query := `
			SELECT SUM(amount) FROM ahead
			WHERE strftime('%%Y-%%m', date) = strftime('%%Y-%%m', now)
		`

		var total sql.NullFloat64
		err := DB.QueryRow(query).Scan(&total)
		if err != nil {
			return 0, err
		}

		if total.Valid {
			return total.Float64, nil
		}
		return 0, nil
	}

	days_left := func() int {
		today := time.Now()
		year, month := today.Year(), today.Month()

		first_of_next_month := time.Date(year, month+1, 1, 0, 0, 0, 0, today.Location())
		days_left := int(first_of_next_month.Sub(today).Hours() / 24)

		return days_left
	}

	var (
		daily_total     float64
		daily_mean      float64
		ahead_total     float64
		overshoot_total float64
		income_total    float64
		daily_forecast  float64
		err             error
	)

	daily_total, err = get_total("spend")
	assert_error("FORECAST error scanning TOTAL of 'spend' table", err)

	daily_mean, err = get_mean()
	assert_error("FORECAST error scanning MEAN of 'spend' table", err)

	ahead_total, err = get_total("ahead")
	assert_error("FORECAST error scanning TOTAL of 'ahead' table", err)

	income_total, err = get_total("income")
	assert_error("FORECAST error scanning TOTAL of 'income' table", err)

	daily_forecast = (income_total - (daily_total + ahead_total)) / (float64(days_left()))

	overshoot_total = income_total - ((daily_total + ahead_total) + (daily_mean * (float64(days_left()))))

	fmt.Printf(`
		Daily Total: %.2f\n
		Ahead Total: %.2f\n
		Income Total: %.2f\n
		Daily Forecast: %.2f\n
		Daily Mean: %.2f\n
		Overshoot: %.2f\n
	`, daily_total, ahead_total, income_total, daily_forecast, daily_mean, overshoot_total)

}

// Get Date from time.Time structure
func UnparseDate(time_struct string) string {

	t, err := time.Parse("2006-01-02 15:04:05-07:00", time_struct)
	assert_error(fmt.Sprintf("Error Parsing time.Time %s", time_struct), err)

	date := t.Format("1-2-2006")

	return date

}

// Date formatting
func ParseDate(unparsed string) time.Time {
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

// Get tag name by id
func GetTagValue(target_id string) string {

	get_tag_name := DB.QueryRow("SELECT name FROM tags WHERE id=?", target_id)

	var result string
	err := get_tag_name.Scan(&result)
	assert_error("Error scanning Get Tag Name by ID", err)

	return result

}

// Inserting or getting a tag
func GetTagID(tag_name string) int64 {

	var tag_id int64
	err := DB.QueryRow("SELECT id FROM tags WHERE name = ?", tag_name).Scan(&tag_id)

	if err == sql.ErrNoRows {
		statement, err := DB.Prepare("INSERT INTO tags(name) VALUES (?)")
		assert_error(fmt.Sprintf("Error 'preparing for insert' new tag(%s) in tags table", tag_name), err)
		defer statement.Close()

		result, err := statement.Exec(tag_name)
		assert_error(fmt.Sprintf("Error 'inserting' new tag(%s) in tags table", tag_name), err)

		tag_id, err := result.LastInsertId()
		assert_error(fmt.Sprintf("Error getting LastInsertId after inserting new tag(%s)", tag_name), err)

		return tag_id
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
func Validate(args []string) (string, int64) {
	action := args[0]

	switch action {
	case "-cd":
		return Create(args[1:])
	default:
		assert_error(fmt.Sprintf("Invalid action '%s'", action), errors.New("Action not recognized"))
	}

	return "There should be an error", 0
}

func main() {
	// fmt.Println(os.Args[1:])
	InitDB()
	fmt.Println(Validate(os.Args[1:]))
	// fmt.Println(Read(1))
}
