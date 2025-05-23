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
func TestCreateDaily(t *testing.T) {
	InitDB()
	var got string

	spend = []string{"test item", "60", "11-30-2001", "testing"}

	got, inserted_id = CreateDaily(spend)
	want := fmt.Sprintf("Daily Spend Created: %s with id %d\n", strings.Join(spend, " "), inserted_id)
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

// Editing
func TestEditDaily(t *testing.T) {
	target_info := 0
	replace_with := "jeep"

	got, replaced_info := EditDaily(inserted_id, target_info, replace_with)
	want := fmt.Sprintf("Edited Daily Spend: %d from %s into %s", inserted_id, replace_with, replaced_info)
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

// Removing
func TestRemoveDaily(t *testing.T) {
	remove := inserted_id
	spend[0] = "jeep"

	got := RemoveDaily(remove)
	want := fmt.Sprintf("Removed Daily Spend: %d %s", inserted_id, strings.Join(spend, " "))
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
