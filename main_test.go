package main

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	var datetests = []struct {
		rawDate       time.Time
		formattedDate string
	}{
		{time.Date(2016, 1, 3, 1, 1, 1, 1, time.UTC), "Jan-03"},
		{time.Date(2016, 1, 15, 1, 1, 1, 1, time.UTC), "Jan-15"},
		{time.Date(2016, 1, 30, 1, 1, 1, 1, time.UTC), "Jan-30"},
	}

	for _, date := range datetests {
		cd := FormatDate(date.rawDate)
		if cd != date.formattedDate {
			t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", date.formattedDate, cd)
		}
	}
}

func TestFormatTime(t *testing.T) {
	var timetests = []struct {
		rawTime       time.Time
		formattedTime string
	}{
		{time.Date(2016, 1, 1, 6, 33, 1, 1, time.UTC), "6:33 a.m."},
		{time.Date(2016, 1, 1, 16, 29, 1, 1, time.UTC), "4:29 p.m."},
		{time.Date(2016, 1, 1, 0, 0, 1, 1, time.UTC), "12:00 a.m."},
	}

	for _, timeTest := range timetests {
		ct := FormatTime(timeTest.rawTime)
		if ct != timeTest.formattedTime {
			t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", timeTest.formattedTime, ct)
		}
	}
}

func TestFindDayWorks(t *testing.T) {
	expectedDay := Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."}
	testFile := "times.example.csv"

	day, _ := findDay(testFile, "Jan-01")

	if day != expectedDay {
		t.Fatalf("Expected %v but was %v", expectedDay, day)
	}
}

func TestFindDayCantFindFile(t *testing.T) {
	_, err := findDay("nofilehere", "Jan-01")

	if err == nil {
		t.Error("Expected a missing file error")
	}
}

func TestFindDayCantFindDay(t *testing.T) {
	testFile := "times.example.csv"
	_, err := findDay(testFile, "Aug-01")

	if err.Error() != "Could not find entry for 'Aug-01' in csv" {
		t.Error("Expected a can't find entry error")
	}
}
