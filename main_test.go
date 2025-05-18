package main

import (
	"testing"
)

func TestInitDB() {
}

func TestCreateDaily(t *testing.T) {
	spend := [4]string{"", "", "", ""}
	got := CreateDaily(spend)
	want := [4]string{"", "", "", ""}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
