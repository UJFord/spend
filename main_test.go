package main

import (
	// "fmt"
	// "strings"
	"testing"
	"time"
)

var (
	daily_inserted_id, monthly_inserted_id, ahead_inserted_id int64
)
var date = time.Date(2011, 11, 30, 0, 0, 0, 0, time.Local)

// Daily
func TestDailyCreate(t *testing.T) {

	InitDB()

	create_tests := []struct {
		name     string
		spend    Daily
		expected string
	}{
		{name: "daily",
			spend: Daily{0, "daily item", 60.0, date, Tag{name: "testing"}, true},
		},
		{name: "monthly",
			spend: Daily{0, "monthly item", 123.0, date, Tag{name: "testing"}, false},
		},
	}

	for _, tt := range create_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Create()

			tt.spend.tag.id = got.tag.id

			if got.isDaily {
				daily_inserted_id = got.id
				tt.spend.id = daily_inserted_id
			} else {

				monthly_inserted_id = got.id
				tt.spend.id = monthly_inserted_id
			}

			if err != nil {
				t.Error(err)
			}

			want := tt.spend

			if got != want {
				t.Errorf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestDailyRead(t *testing.T) {

	read_tests := []struct {
		name  string
		spend Daily
	}{
		{name: "daily",
			spend: Daily{daily_inserted_id, "daily item", 60.0, date, Tag{name: "testing"}, true},
		},
		{name: "monthly",
			spend: Daily{monthly_inserted_id, "monthly item", 123.0, date, Tag{name: "testing"}, false},
		},
	}

	for _, tt := range read_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Read(tt.spend.id)
			if err != nil {
				t.Error(err)
			}

			tt.spend.tag.id = got.tag.id

			want := tt.spend

			if got != tt.spend {
				t.Errorf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestDailyEdit(t *testing.T) {

	edit_tests := []struct {
		name            string
		value           any
		field           int
		new_value_spend Daily
	}{
		{name: "daily",
			field:           0,
			value:           "daily replacement",
			new_value_spend: Daily{daily_inserted_id, "daily replacement", 60.0, date, Tag{name: "testing"}, true},
		},
		{name: "tag replace",
			field:           3,
			value:           "tag replacement",
			new_value_spend: Daily{monthly_inserted_id, "monthly item", 123.0, date, Tag{name: "tag replacement"}, false},
		},
		{name: "monthly",
			field:           1,
			value:           456.7,
			new_value_spend: Daily{monthly_inserted_id, "monthly item", 456.7, date, Tag{name: "tag replacement"}, false},
		},
	}

	for _, tt := range edit_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.new_value_spend.Edit(tt.field, tt.value)
			if err != nil {
				t.Error(err)
			}

			tt.new_value_spend.tag.id = got.tag.id

			want := tt.new_value_spend

			if got != want {
				t.Errorf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestDailyRemove(t *testing.T) {

	remove_test := []struct {
		name  string
		spend Daily
	}{
		{name: "daily",
			spend: Daily{daily_inserted_id, "daily replacement", 60.0, date, Tag{name: "testing"}, true},
		},
		{name: "monthly",
			spend: Daily{monthly_inserted_id, "monthly item", 456.7, date, Tag{name: "tag replacement"}, false},
		},
	}

	for _, tt := range remove_test {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Remove()
			if err != nil {
				t.Error(err)
			}

			tt.spend.tag.id = got.tag.id

			want := tt.spend

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

// Spending Ahead
func TestAheadCreate(t *testing.T) {

	create_tests := []struct {
		name  string
		spend Ahead
	}{
		{name: "create",
			spend: Ahead{0, 999.0, date},
		},
	}

	for _, tt := range create_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Create()
			if err != nil {
				t.Error(err)
			}

			ahead_inserted_id = got.id
			tt.spend.id = ahead_inserted_id

			want := tt.spend

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

func TestAheadRead(t *testing.T) {

	read_tests := []struct {
		name  string
		spend Ahead
	}{
		{name: "read",
			spend: Ahead{ahead_inserted_id, 999.0, date},
		},
	}

	for _, tt := range read_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Read(ahead_inserted_id)
			if err != nil {
				t.Error(err)
			}

			want := tt.spend

			if got != tt.spend {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

func TestAheadRemove(t *testing.T) {

	remove_tests := []struct {
		name  string
		spend Ahead
	}{
		{name: "remove",
			spend: Ahead{ahead_inserted_id, 999.0, date},
		},
	}

	for _, tt := range remove_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Read(tt.spend.id)
			if err != nil {
				t.Error(err)
			}

			want := tt.spend

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

// Tags
func TestTagEdit(t *testing.T) {

	edit_tests := []struct {
		name    string
		tag     Tag
		replace string
		new_tag Tag
	}{
		{name: "editing",
			tag:     Tag{0, "testing"},
			replace: "testing tag edit",
			new_tag: Tag{0, "testing tag edit"},
		},
		{name: "unnaming",
			tag:     Tag{0, "testing tag edit"},
			replace: "",
			new_tag: Tag{0, ""},
		},
	}

	for _, tt := range edit_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.tag.Edit(tt.replace)
			if err != nil {
				t.Error(err)
			}

			tt.new_tag.id = got.id

			want := tt.new_tag

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

// Income
func TestIncomeCreate(t *testing.T) {

	create_tests := []struct {
		name   string
		income Income
	}{
		{name: "10K",
			income: Income{0, 10000.0, date},
		},
	}

	for _, tt := range create_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.income.Create()
			if err != nil {
				t.Error(err)
			}

			tt.income.id = got.id

			want := tt.income

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}
