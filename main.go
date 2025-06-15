package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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

type Daily struct {
	id      int64
	item    string
	amount  float64
	date    time.Time
	tag     Tag
	isDaily bool
}

type DailyActions interface {
	Create() (Daily, error)
	Read(int64) (Daily, error)
	Edit(int, any) (Daily, error)
	Remove() (Daily, error)
}

type Ahead struct {
	id     int64
	amount float64
	date   time.Time
}

type AheadActions interface {
	Create() (Ahead, error)
	Read(int64) (Ahead, error)
	Remove() (Ahead, error)
}

type Tag struct {
	id   int64
	name string
}

type TagActions interface {
	SetID() (Tag, error)
	SetValue() (Tag, error)

	Edit(string) (Tag, error)
}

type Income struct {
	id     int64
	amount float64
	date   time.Time
}

type IncomeActions interface {
	Create() (Income, error)
	Read() (Income, error)
	Edit() (Income, error)
	Remove() (Income, error)
}

type Forecast struct {
	daily_total     float64
	ahead_total     float64
	income_total    float64
	daily_forecast  float64
	daily_mean      float64
	overshoot_total float64
}

type ForecastActions interface {
	Update() (Forecast, error)
}

// Initalize db
func InitDB() error {
	target_db_file = "./spend.db"
	var err error

	DB, err = sql.Open("sqlite3", target_db_file)
	if err != nil {
		return fmt.Errorf("error initializing db: '%w'", err)
	}

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
		name TEXT DEFAULT 'unnamed'
	);
	CREATE TABLE IF NOT EXISTS ahead(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		date TEXT DEFAULT '1970-01-01 00:00:00+00:00'
	);
	CREATE TABLE IF NOT EXISTS income(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		date TEXT DEFAULT '1970-01-01 00:00:00+00:00'
	);
	`
	_, err = DB.Exec(init_schema)
	if err != nil {
		return fmt.Errorf("error executing init db: '%w'", err)
	}

	return nil
}

// Create
func (s Daily) Create() (Daily, error) {

	var err error
	s.tag, err = s.tag.SetID()
	if err != nil {
		return Daily{}, fmt.Errorf("create error getting tag id: '%w'", err)
	}

	insert_stmt, err := DB.Prepare(fmt.Sprintf(`
		INSERT INTO spend(item, amount, date, tag_id, is_daily)
		VALUES (?, ?, ?, ?, ?)
	`))
	if err != nil {
		return Daily{}, fmt.Errorf("create error preparing insert statement: '%w'", err)
	}
	defer insert_stmt.Close()

	exec_insert_stmt, err := insert_stmt.Exec(s.item, s.amount, s.date, s.tag.id, s.isDaily)
	if err != nil {
		return Daily{}, fmt.Errorf("create error executing insert statement: '%w'", err)
	}

	s.id, err = exec_insert_stmt.LastInsertId()
	if err != nil {
		return Daily{}, fmt.Errorf("create error fetching last insert id: '%w'", err)
	}

	return s, nil
}

func (a Ahead) Create() (Ahead, error) {

	insert_stmt, err := DB.Prepare("INSERT INTO ahead(amount, date) VALUES(?, ?)")

	if err != nil {
		return Ahead{}, fmt.Errorf("create ahead error preparing insert statement: '%w'", err)
	}
	defer insert_stmt.Close()

	exec_insert_stmt, err := insert_stmt.Exec(a.amount, a.date)
	if err != nil {
		return Ahead{}, fmt.Errorf("create ahead error executing insert statement: '%w'", err)
	}

	id_of_inserted, err := exec_insert_stmt.LastInsertId()
	if err != nil {
		return Ahead{}, fmt.Errorf("create ahead error fetching last insert id: '%w'", err)
	}

	a.id = id_of_inserted

	return a, nil
}

// Get info by id
func (s Daily) Read(target_id int64) (Daily, error) {

	get := DB.QueryRow("SELECT item, amount, date, tag_id, is_daily FROM spend WHERE id=?", target_id)

	var unparsed_date string
	err := get.Scan(&s.item, &s.amount, &unparsed_date, &s.tag.id, &s.isDaily)
	if err != nil {
		return Daily{}, fmt.Errorf("read error scanning read query: '%w'", err)
	}

	s.id = target_id

	layout := "2006-01-02 15:04:05-07:00"
	s.date, err = time.Parse(layout, unparsed_date)
	if err != nil {
		return Daily{}, fmt.Errorf("read error parsing date from db: '%w'", err)
	}

	s.tag, err = s.tag.SetValue()
	if err != nil {
		return Daily{}, err
	}

	return s, nil
}

func (a Ahead) Read(target_id int64) (Ahead, error) {
	get := DB.QueryRow("SELECT amount, date FROM ahead WHERE id=?", target_id)

	var unparsed_date string
	err := get.Scan(&a.amount, &unparsed_date)
	if err != nil {
		return Ahead{}, fmt.Errorf("read ahead error scanning query stetement: '%w'", err)
	}

	a.id = target_id

	layout := "2006-01-02 15:04:05-07:00"
	a.date, err = time.Parse(layout, unparsed_date)
	if err != nil {
		return Ahead{}, fmt.Errorf("read error parsing date from db: '%w'", err)
	}

	return a, nil
}

// Edit a spend
func (s Daily) Edit(target_field int, replace_value any) (Daily, error) {

	var err error

	if target_field < 0 || target_field > 4 {
		log.Fatalf("EDIT error choosing target info: Out of Bounds")
	}

	var target string
	switch target_field {
	case 0:
		target = "item"
	case 1:
		target = "amount"
	case 2:
		target = "date"
	case 3:
		target = "tag_id"

		if _, ok := replace_value.(string); ok {

			s.tag, err = s.tag.SetID()
			if err != nil {
				return Daily{}, err
			}

			replace_value = strconv.FormatInt(s.tag.id, 10)
		}
	case 4:
		target = "is_daily"
	}

	edit_stmt, err := DB.Prepare(fmt.Sprintf("UPDATE spend SET %s = ? WHERE id = ?", target))
	if err != nil {
		return Daily{}, fmt.Errorf("edit error preparing edit statement: '%w'", err)
	}
	defer edit_stmt.Close()

	_, err = edit_stmt.Exec(replace_value, s.id)
	if err != nil {
		return Daily{}, fmt.Errorf("edit error executing edit statement: '%w'", err)
	}

	s, err = s.Read(s.id)
	if err != nil {
		return Daily{}, err
	}

	return s, nil
}

// Remove spend
func (s Daily) Remove() (Daily, error) {

	s, err := s.Read(s.id)
	if err != nil {
		return Daily{}, err
	}

	delete_stmt, err := DB.Prepare("DELETE FROM spend WHERE id=?")
	if err != nil {
		return Daily{}, fmt.Errorf("remove error preparing delete statement: '%w'", err)
	}

	_, err = delete_stmt.Exec(s.id)
	if err != nil {
		return Daily{}, fmt.Errorf("remove error executing delete statement: '%w'", err)
	}

	return s, nil
}

func (a Ahead) Remove() (Ahead, error) {

	a, err := a.Read(a.id)
	if err != nil {
		return Ahead{}, err
	}

	delete_stmt, err := DB.Prepare("DELETE FROM ahead WHERE id=?")
	if err != nil {
		return Ahead{}, fmt.Errorf("remove error preparing delete statement: '%w'", err)
	}

	_, err = delete_stmt.Exec(a.id)
	if err != nil {
		return Ahead{}, fmt.Errorf("remove error executing delete statement: '%w'", err)
	}

	return a, nil
}

// Forecast
func (f Forecast) Update() (Forecast, error) {
	// ahead sum
	get_total := func(table string) (float64, error) {

		query := fmt.Sprintf(`
			SELECT SUM(amount) FROM %s
			WHERE strftime('%%Y-%%m', date) = strftime('%%Y-%%m', 'now')
		`, table)

		var total sql.NullFloat64
		err := DB.QueryRow(query).Scan(&total)
		if err != nil {
			return 0, err
		}

		if !total.Valid {
			return 0, nil
		}

		return total.Float64, nil
	}

	get_mean := func() (float64, error) {
		query := `
			SELECT SUM(amount) FROM ahead
			WHERE strftime('%%Y-%%m', date) = strftime('%%Y-%%m', 'now')
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

	var err error

	f.daily_total, err = get_total("spend")
	if err != nil {
		return Forecast{}, fmt.Errorf("forecast error fetching daily total: '%w'", err)
	}

	f.daily_mean, err = get_mean()
	if err != nil {
		return Forecast{}, fmt.Errorf("forecast error fetching daily mean: '%w'", err)
	}

	f.ahead_total, err = get_total("ahead")
	if err != nil {
		return Forecast{}, fmt.Errorf("forecast error fetching ahead mean: '%w'", err)
	}

	f.income_total, err = get_total("income")
	if err != nil {
		return Forecast{}, fmt.Errorf("forecast error fetching income total: '%w'", err)
	}

	f.daily_forecast = (f.income_total - (f.daily_total + f.ahead_total)) / (float64(days_left()))

	f.overshoot_total = f.income_total - ((f.daily_total + f.ahead_total) + (f.daily_mean * (float64(days_left()))))

	return f, nil
}

// DATE
func UnparseDate(time_struct string) (string, error) {

	t, err := time.Parse("2006-01-02 15:04:05-07:00", time_struct)
	if err != nil {
		return "", fmt.Errorf("unparse date error parsing date: '%w'", err)
	}

	date := t.Format("1-2-2006")

	return date, nil

}

func ParseDate(unparsed string) (time.Time, error) {
	split_unparsed := strings.Split(unparsed, "-")

	month, err := strconv.Atoi(split_unparsed[0])
	if err != nil {
		return time.Time{}, errors.New("error converting month to int")
	}

	day, err := strconv.Atoi(split_unparsed[1])
	if err != nil {
		return time.Time{}, errors.New("error converting day to int")
	}

	year, err := strconv.Atoi(split_unparsed[2])
	if err != nil {
		return time.Time{}, errors.New("error converting year to int")
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if !(t.Year() == year && int(t.Month()) == month && t.Day() == day) {
		return time.Time{}, errors.New("ivalid date format")
	}

	return t, nil
}

// TAG
func (t Tag) Edit(replace string) (Tag, error) {

	var err error

	t, err = t.SetID()
	if err != nil {
		return Tag{}, err
	}

	edit_stmt, err := DB.Prepare("UPDATE tags SET name = ? WHERE id = ?")
	if err != nil {
		return Tag{}, fmt.Errorf("tag edit error preparing update statement: '%w'", err)
	}
	defer edit_stmt.Close()

	_, err = edit_stmt.Exec(replace, t.id)
	if err != nil {
		return Tag{}, fmt.Errorf("tag edit error executing edit statement: '%w'", err)
	}

	t, err = t.SetValue()
	if err != nil {
		return Tag{}, err
	}

	return t, nil
}

func (t Tag) SetValue() (Tag, error) {

	get_tag_name := DB.QueryRow("SELECT name FROM tags WHERE id=?", t.id)

	err := get_tag_name.Scan(&t.name)
	if err != nil {
		return Tag{}, fmt.Errorf("tag value error scanning get tag name statement: '%w'", err)
	}

	return t, nil

}

func (t Tag) SetID() (Tag, error) {

	err := DB.QueryRow("SELECT id FROM tags WHERE name = ?", t.name).Scan(&t.id)

	if err == sql.ErrNoRows {

		statement, err := DB.Prepare("INSERT INTO tags(name) VALUES (?)")
		if err != nil {
			return Tag{}, fmt.Errorf("tag error preparing insert statement: '%w'", err)
		}
		defer statement.Close()

		result, err := statement.Exec(t.name)
		if err != nil {
			return Tag{}, fmt.Errorf("tag error executing insert statement: '%w'", err)
		}

		t.id, err = result.LastInsertId()
		if err != nil {
			return Tag{}, fmt.Errorf("tag error fetching last insert id: '%w'", err)
		}

	} else if err != nil {

		return Tag{}, fmt.Errorf("tag error query tag name: '%w'", err)

	}

	return t, nil

}

// INCOME
func (i Income) Create() (Income, error) {

	create_stmt, err := DB.Prepare(`
		INSERT INTO income(amount, date)
		VALUES(?, ?)
	`)
	if err != nil {
		return Income{}, fmt.Errorf("income add error preparing create statement: '%w'", err)
	}
	defer create_stmt.Close()

	exec_create_stmt, err := create_stmt.Exec(i.amount, i.date)
	if err != nil {
		return Income{}, fmt.Errorf("income add error executing create statement: '%w'", err)
	}

	i.id, err = exec_create_stmt.LastInsertId()
	if err != nil {
		return Income{}, fmt.Errorf("income add error fetching last insert id: '%w'", err)
	}

	return i, nil
}

func (i Income) Read() (Income, error) {

	get := DB.QueryRow("SELECT amount, date FROM income WHERE id=?", i.id)

	var unparsed_date string
	err := get.Scan(&i.amount, &unparsed_date)
	if err != nil {
		return Income{}, fmt.Errorf("read ahead error scanning query stetement: '%w'", err)
	}

	layout := "2006-01-02 15:04:05-07:00"
	i.date, err = time.Parse(layout, unparsed_date)
	if err != nil {
		return Income{}, fmt.Errorf("read error parsing date from db: '%w'", err)
	}

	return i, nil
}

func main() {

	if err := InitDB(); err != nil {
		log.Fatal(err)
	}

	f := Forecast{}
	if f, err := f.Update(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(f)
	}
}
