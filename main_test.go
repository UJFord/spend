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
func TestCreate(t *testing.T) {

	InitDB()

	create_tests := []struct {
		name     string
		spend    Spend
		expected string
	}{
		{name: "daily",
			spend:    Daily{0, "daily item", 60.0, date, "testing", true},
			expected: "CREATE: 'daily item 60 11-30-2001 testing'",
		},
		{name: "monthly",
			spend:    Daily{0, "monthly item", 123.0, date, "testing", false},
			expected: "CREATE: 'monthly item 60 11-30-2001 testing'",
		},
	}

	for _, tt := range create_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Create()

			if got.isDaily {
				daily_inserted_id = got.id
				tt.spend = tt.spend.SetID(daily_inserted_id)
			} else {
				monthly_inserted_id = got.id
				tt.spend = tt.spend.SetID(monthly_inserted_id)
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

func TestRead(t *testing.T) {

	read_tests := []struct {
		name  string
		spend Spend
	}{
		{name: "daily",
			spend: Daily{daily_inserted_id, "daily item", 60.0, date, "testing", true},
		},
		{name: "monthly",
			spend: Daily{monthly_inserted_id, "monthly item", 123.0, date, "testing", false},
		},
	}

	for _, tt := range read_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Read(tt.spend.GetStruct().id)
			if err != nil {
				t.Error(err)
			}

			want := tt.spend

			if got != tt.spend {
				t.Errorf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestEdit(t *testing.T) {

	edit_tests := []struct {
		name            string
		value           any
		field           int
		new_value_spend Spend
	}{
		{name: "daily",
			field:           0,
			value:           "daily replacement",
			new_value_spend: Daily{daily_inserted_id, "daily replacement", 60.0, date, "testing", true},
		},
		{name: "tag replace",
			field:           3,
			value:           "tag replacement",
			new_value_spend: Daily{monthly_inserted_id, "monthly item", 123.0, date, "tag replacement", false},
		},
		{name: "monthly",
			field:           1,
			value:           456.7,
			new_value_spend: Daily{monthly_inserted_id, "monthly item", 456.7, date, "tag replacement", false},
		},
	}

	for _, tt := range edit_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.new_value_spend.Edit(tt.field, tt.value)

			if err != nil {
				t.Error(err)
			}

			want := tt.new_value_spend

			if got != want {
				t.Errorf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestRemove(t *testing.T) {

	remove_test := []struct {
		name  string
		spend Daily
	}{
		{name: "daily",
			spend: Daily{daily_inserted_id, "daily replacement", 60.0, date, "testing", true},
		},
		{name: "monthly",
			spend: Daily{monthly_inserted_id, "monthly item", 456.7, date, "tag replacement", false},
		},
	}

	for _, tt := range remove_test {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.spend.Remove()

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

// Spending Ahead
func TestCreateAhead(t *testing.T) {

	create_tests := []struct {
		name  string
		spend SpendAhead
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
			tt.spend = tt.spend.SetID(got.id)

			want := tt.spend

			if got != want {
				t.Errorf("got '%#v' want '%#v'", got, want)
			}
		})
	}
}

func TestReadAhead(t *testing.T) {

	read_tests := []struct {
		name  string
		spend SpendAhead
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

func TestRemoveAhead(t *testing.T) {

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
		name     string
		existing string
		replace  string
	}{
		{name: "editing",
			existing: "testing",
			replace:  "testing tag edit",
		},
		{name: "unnaming",
			existing: "testing tag edit",
			replace:  "",
		},
	}

	for _, tt := range edit_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := TagEdit(tt.existing, tt.replace)
			if err != nil {
				t.Error(err)
			}

			want := tt.replace

			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
	}
}
