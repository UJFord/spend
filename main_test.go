package main

import (
	"testing"
)

func TestInitDB() {
}

func TestCreate(t *testing.T) {
	spend := [5]string{"", "", "", "", ""}
	got := Create(spend)
	want := [5]string{"", "", "", "", ""}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
