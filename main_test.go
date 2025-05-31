package main

import (
	"fmt"
	"strings"
	"testing"
)

var (
	inserted_id, ahead_inserted_id int64
)

// Adding
func TestCreate(t *testing.T) {

	InitDB()

	var got string

	t.Run("empty frequency", func(t *testing.T) {
		spend := []string{"daily item", "60", "11-30-2001", "testing", ""}

		got, inserted_id = Create(spend)
		want := fmt.Sprintf("CREATE spend created: %s with id %d\n",
			strings.Join(spend, " "), inserted_id)

		log_error(t, got, want)
	})

	t.Run("daily", func(t *testing.T) {
		spend := []string{"daily item", "60", "11-30-2001", "testing", "daily"}

		got, inserted_id = Create(spend)
		want := fmt.Sprintf("CREATE spend created: %s with id %d\n",
			strings.Join(spend, " "), inserted_id)

		log_error(t, got, want)
	})

	t.Run("monthly", func(t *testing.T) {
		spend := []string{"monthly item", "4500", "11-30-2001", "testing", "monthly"}

		got, inserted_id = Create(spend)
		want := fmt.Sprintf("CREATE spend created: %s with id %d\n",
			strings.Join(spend, " "), inserted_id)

		log_error(t, got, want)
	})
}

// // Reading
// func TestRead(t *testing.T) {
//
//		t.Run("daily", func(t *testing.T) {
//			_, got := Read(inserted_id, d)
//
//			want := fmt.Sprintf("READ %s info: %d %s",
//				d, inserted_id, strings.Join(spend, " "))
//			log_error(t, got, want)
//		})
//
//		t.Run("monthly", func(t *testing.T) {
//			_, got := Read(inserted_id, m)
//
//			want := fmt.Sprintf("READ %s info: %d %s",
//				m, inserted_id, strings.Join(spend, " "))
//			log_error(t, got, want)
//		})
//	}
//
// // Editing
// func TestEdit(t *testing.T) {
//
//		t.Run("daily", func(t *testing.T) {
//			target_info := 0
//			replace_with := "jeep"
//
//			got, replaced_info := Edit(inserted_id, target_info, replace_with, d)
//			want := fmt.Sprintf("EDIT edited %s spend: id(%d) from (%s) into (%s)",
//				d, inserted_id, replace_with, replaced_info)
//
//			log_error(t, got, want)
//		})
//
//		t.Run("monthly", func(t *testing.T) {
//			target_info := 0
//			replace_with := "jeep"
//
//			got, replaced_info := Edit(inserted_id, target_info, replace_with, m)
//			want := fmt.Sprintf("EDIT edited %s spend: id(%d) from (%s) into (%s)",
//				m, inserted_id, replace_with, replaced_info)
//
//			log_error(t, got, want)
//		})
//	}
//
// // Removing
// func TestRemove(t *testing.T) {
//
//		t.Run("daily", func(t *testing.T) {
//			spend[0] = "jeep"
//
//			got := Remove(inserted_id, d)
//			want := fmt.Sprintf("REMOVE removed %s spend: %d %s",
//				d, inserted_id, strings.Join(spend, " "))
//
//			log_error(t, got, want)
//		})
//
//		t.Run("monthly", func(t *testing.T) {
//			spend[0] = "jeep"
//
//			got := Remove(inserted_id, m)
//			want := fmt.Sprintf("REMOVE removed %s spend: %d %s",
//				m, inserted_id, strings.Join(spend, " "))
//
//			log_error(t, got, want)
//		})
//	}
//
// // Spending Ahead
// var (
//
//	ahead_amount = 999.00
//	ahead_date   = "1-20-2025"
//
// )
//
// func TestCreateAhead(t *testing.T) {
//
//		var got string
//		got, ahead_inserted_id = CreateAhead(ahead_amount, ahead_date)
//		want := fmt.Sprintf("CREATE AHEAD spending amount(%.2f) date(%s) ahead with id(%d)",
//			ahead_amount,
//			ahead_date,
//			ahead_inserted_id)
//
//		log_error(t, got, want)
//	}
//
// func TestReadAhead(t *testing.T) {
//
//		_, got := ReadAhead(ahead_inserted_id)
//		want := fmt.Sprintf("READ AHEAD id(%d) amount(%.2f) date(%s)", ahead_inserted_id, ahead_amount, ahead_date)
//
//		log_error(t, got, want)
//	}
//
// func TestRemoveAhead(t *testing.T) {
//
//		got := RemoveAhead(ahead_inserted_id)
//		want := fmt.Sprintf("REMOVE AHEAD spending amount(%.2f) ahead with id(%d)", ahead_amount, ahead_inserted_id)
//
//		log_error(t, got, want)
//	}
//
// return error
func log_error(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
