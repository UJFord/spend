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
		want := fmt.Sprintf("Daily Spend Created: %s with id %d\n", strings.Join(spend[1:], " "), inserted_id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	// Removing
	// 	t.Run("removing created manok", func(t *testing.T) {
	// 		remove := []string{"-rd", strconv.FormatInt(inserted_id, 10)}
	//
	// 		got, _ := Validate(remove)
	// 	})
}
