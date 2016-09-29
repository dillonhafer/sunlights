package main

import (
	"testing"
	"time"
)

func assertEqual(t *testing.T, given Day, expected Day) {
	if given != expected {
		t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", expected, given)
	}
}

func TestCurrentDate(t *testing.T) {
	var datetests = []struct {
		in  time.Time
		out string
	}{
		{time.Date(2016, 1, 3, 1, 1, 1, 1, time.UTC), "Jan-03"},
		{time.Date(2016, 1, 15, 1, 1, 1, 1, time.UTC), "Jan-15"},
		{time.Date(2016, 1, 30, 1, 1, 1, 1, time.UTC), "Jan-30"},
	}

	for _, rawDate := range datetests {
		date := CurrentDate(rawDate.in)
		if date != rawDate.out {
			t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", rawDate.out, date)
		}
	}
}

func TestFindDayWorks(t *testing.T) {
	expectedDay := Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."}
	testFile := "times.example.csv"

	day, _ := findDay(testFile, "Jan-01")

	assertEqual(t, day, expectedDay)
}

func TestFindDayCantFindFile(t *testing.T) {
	_, err := findDay("nofilehere", "Jan-01")

	if err == nil {
		t.Error("Expected an error")
	}
}
