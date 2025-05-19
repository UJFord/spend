package main

import (
	"testing"
)

func TestCreateDaily(t *testing.T) {
	spend := []string{"-cd", "manok", "60", "1-1-2001", "ulam"}
	got := validate_input(spend)
	want := "Created daily spend: manok 60 ulam"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
