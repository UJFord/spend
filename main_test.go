package main

import (
	"fmt"
	"strings"
	"testing"
)

var (
	inserted_id int64
	spend       []string
)

// Adding
func TestCreate(t *testing.T) {
	InitDB()
	var got string

	t.Run("daily", func(t *testing.T) {
		spend = []string{"test item", "60", "11-30-2001", "testing"}

		got, inserted_id = Create(spend, "daily")
		want := fmt.Sprintf("daily Spend Created: %s with id %d\n", strings.Join(spend, " "), inserted_id)
		log_error(t, got, want)
	})

	t.Run("monthly", func(t *testing.T) {
		spend = []string{"rent", "4500", "11-30-2001", "rent"}

		got, inserted_id = Create(spend, "monthly")
		want := fmt.Sprintf("monthly Spend Created: %s with id %d\n", strings.Join(spend, " "), inserted_id)
		log_error(t, got, want)
	})
}

// Reading
func TestReadDaily(t *testing.T) {

	_, daily_info := ReadDaily(inserted_id)

	got := daily_info
	want := fmt.Sprintf("Daily info: %d %s", inserted_id, strings.Join(spend, " "))
	log_error(t, got, want)
}

// Editing
func TestEditDaily(t *testing.T) {
	t.Helper()
	target_info := 0
	replace_with := "jeep"

	got, replaced_info := EditDaily(inserted_id, target_info, replace_with)
	want := fmt.Sprintf("Edited Daily Spend: id(%d) from (%s) into (%s)", inserted_id, replace_with, replaced_info)
	log_error(t, got, want)
}

// Removing
func TestRemoveDaily(t *testing.T) {
	t.Helper()
	remove := inserted_id
	spend[0] = "jeep"

	got := RemoveDaily(remove)
	want := fmt.Sprintf("Removed Daily Spend: %d %s", inserted_id, strings.Join(spend, " "))
	log_error(t, got, want)
}

// return error
func log_error(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
