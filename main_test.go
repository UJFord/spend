package main

import (
	"fmt"
	"strings"
	"testing"
)

var (
	daily_inserted_id, monthly_inserted_id int64
	daily_spend, monthly_spend             []string
	d, m                                   string
)

// Adding
func TestCreate(t *testing.T) {
	InitDB()
	d = "daily"
	m = "monthly"
	var got string

	t.Run(d, func(t *testing.T) {
		daily_spend = []string{"test item", "60", "11-30-2001", "testing"}

		got, daily_inserted_id = Create(daily_spend, d)
		want := fmt.Sprintf("CREATE daily spend created: %s with id %d\n",
			strings.Join(daily_spend, " "), daily_inserted_id)

		log_error(t, got, want)
	})

	t.Run(m, func(t *testing.T) {
		monthly_spend = []string{"rent", "4500", "11-30-2001", "rent"}

		got, monthly_inserted_id = Create(monthly_spend, m)
		want := fmt.Sprintf("CREATE monthly spend created: %s with id %d\n",
			strings.Join(monthly_spend, " "), monthly_inserted_id)

		log_error(t, got, want)
	})
}

// Reading
func TestRead(t *testing.T) {

	t.Run("daily", func(t *testing.T) {
		_, got := Read(daily_inserted_id, d)

		want := fmt.Sprintf("READ %s info: %d %s",
			d, daily_inserted_id, strings.Join(daily_spend, " "))
		log_error(t, got, want)
	})

	t.Run("monthly", func(t *testing.T) {
		_, got := Read(monthly_inserted_id, m)

		want := fmt.Sprintf("READ %s info: %d %s",
			m, monthly_inserted_id, strings.Join(monthly_spend, " "))
		log_error(t, got, want)
	})
}

// Editing
func TestEdit(t *testing.T) {

	t.Run("daily", func(t *testing.T) {
		target_info := 0
		replace_with := "jeep"

		got, replaced_info := Edit(daily_inserted_id, target_info, replace_with, d)
		want := fmt.Sprintf("EDIT edited %s spend: id(%d) from (%s) into (%s)",
			d, daily_inserted_id, replace_with, replaced_info)

		log_error(t, got, want)
	})

	t.Run("monthly", func(t *testing.T) {
		target_info := 0
		replace_with := "jeep"

		got, replaced_info := Edit(monthly_inserted_id, target_info, replace_with, m)
		want := fmt.Sprintf("EDIT edited %s spend: id(%d) from (%s) into (%s)",
			m, monthly_inserted_id, replace_with, replaced_info)

		log_error(t, got, want)
	})
}

// Removing
func TestRemove(t *testing.T) {

	t.Run("daily", func(t *testing.T) {
		daily_spend[0] = "jeep"

		got := Remove(daily_inserted_id, d)
		want := fmt.Sprintf("REMOVE removed %s spend: %d %s",
			d, daily_inserted_id, strings.Join(daily_spend, " "))

		log_error(t, got, want)
	})

	t.Run("monthly", func(t *testing.T) {
		monthly_spend[0] = "jeep"

		got := Remove(monthly_inserted_id, m)
		want := fmt.Sprintf("REMOVE removed %s spend: %d %s",
			m, monthly_inserted_id, strings.Join(monthly_spend, " "))

		log_error(t, got, want)
	})
}

// return error
func log_error(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
