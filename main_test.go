package main

import (
	"testing"
)

func TestInitDB() {
}

func TestCreateDaily(t *testing.T) {
	spend := []string{"-cd", "manok", "60", "ulam"}
	validate_input(spend)
}
