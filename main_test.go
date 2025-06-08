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

// Adding
func TestCreate(t *testing.T) {

	InitDB()

	create_tests := []struct {
		name     string
		spend    Daily
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

			s, err := tt.spend.Create()

			if s.isDaily {
				daily_inserted_id = s.id
				tt.spend.id = daily_inserted_id
			} else {
				monthly_inserted_id = s.id
				tt.spend.id = monthly_inserted_id
			}

			if err != nil {
				t.Error(err)
			}

			if s != tt.spend {
				t.Errorf("got %#v want %#v", s, tt.spend)
			}
		})
	}
}

// Reading
func TestRead(t *testing.T) {

	read_tests := []struct {
		name  string
		spend Daily
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

			got, err := tt.spend.Read(tt.spend.id)

			if err != nil {
				t.Error(err)
			}

			if got != tt.spend {
				t.Errorf("got %#v want %#v", got, tt.spend)
			}
		})
	}
}

// // Editing
func TestEdit(t *testing.T) {

	read_tests := []struct {
		name            string
		value           any
		field           int
		new_value_spend Daily
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

	for _, tt := range read_tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.new_value_spend.Edit(tt.field, tt.value)

			if err != nil {
				t.Error(err)
			}

			if got != tt.new_value_spend {
				t.Errorf("got %#v want %#v", got, tt.new_value_spend)
			}
		})
	}
}

//
// // Removing
// func TestRemove(t *testing.T) {
//
// 	t.Run("daily", func(t *testing.T) {
// 		spend := []string{"daily item replacement", "60", "11-30-2001", "testing", "daily"}
//
// 		got := Remove(daily_inserted_id)
// 		want := fmt.Sprintf("REMOVE removed spend: %d %s",
// 			daily_inserted_id, strings.Join(spend, " "))
//
// 		log_error(t, got, want)
// 	})
//
// 	t.Run("monthly", func(t *testing.T) {
// 		spend := []string{"monthly item replacement", "4500", "11-30-2001", "testing", "monthly"}
//
// 		got := Remove(monthly_inserted_id)
// 		want := fmt.Sprintf("REMOVE removed spend: %d %s",
// 			monthly_inserted_id, strings.Join(spend, " "))
//
// 		log_error(t, got, want)
// 	})
// }
//
// // Spending Ahead
// var (
// 	ahead_amount = 999.00
// 	ahead_date   = "1-20-2025"
// )
//
// func TestCreateAhead(t *testing.T) {
//
// 	var got string
// 	got, ahead_inserted_id = CreateAhead(ahead_amount, ahead_date)
// 	want := fmt.Sprintf("CREATE AHEAD spending amount(%.2f) date(%s) ahead with id(%d)",
// 		ahead_amount,
// 		ahead_date,
// 		ahead_inserted_id)
//
// 	log_error(t, got, want)
// }
//
// func TestReadAhead(t *testing.T) {
//
// 	_, got := ReadAhead(ahead_inserted_id)
// 	want := fmt.Sprintf("READ AHEAD id(%d) amount(%.2f) date(%s)", ahead_inserted_id, ahead_amount, ahead_date)
//
// 	log_error(t, got, want)
// }
//
// func TestRemoveAhead(t *testing.T) {
//
// 	got := RemoveAhead(ahead_inserted_id)
// 	want := fmt.Sprintf("REMOVE AHEAD spending amount(%.2f) ahead with id(%d)", ahead_amount, ahead_inserted_id)
//
// 	log_error(t, got, want)
// }
