package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var inserted_id int64

func TestValidate(t *testing.T) {
	InitDB()

	// Adding
	t.Run("adding manok", func(t *testing.T) {
		spend := []string{"manok", "60", "1-1-2001", "ulam"}

		got, inserted_id := CreateDaily(spend)
		want := fmt.Sprintf("Daily Spend Created: %s with id %d\n", strings.Join(spend, " "), inserted_id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	// Removing
	t.Run("removing created manok", func(t *testing.T) {
		remove := inserted_id

		got := RemoveDaily(remove)
	})
}
