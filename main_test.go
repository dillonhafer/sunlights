package main

import (
	"testing"
)

func assertEqual(t *testing.T, given Day, expected Day) {
	if given != expected {
		t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", expected, given)
	}
}

func TestFindDay(t *testing.T) {
	expectedDay := Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."}
	testFile := "times.example.csv"

	day := findDay(testFile, "Jan-01")

	assertEqual(t, day, expectedDay)
}
