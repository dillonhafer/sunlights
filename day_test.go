package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFindDayWorks(t *testing.T) {
	expectedDay := Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."}
	days := []Day{expectedDay}

	day, _ := findDay(days, "Jan-01")

	if day != expectedDay {
		t.Fatalf("Expected %v but was %v", expectedDay, day)
	}
}

func TestFindDayCantFindDay(t *testing.T) {
	expectedDay := Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."}
	days := []Day{expectedDay}
	_, err := findDay(days, "Aug-01")

	if err.Error() != "Could not find entry for 'Aug-01' in config.json" {
		t.Error("Expected a can't find entry error")
	}
}

func TestNewToday(t *testing.T) {
	days := []Day{
		Day{day: "Jan-01", sunrise: "7:15 a.m.", sunset: "4:29 p.m."},
		Day{day: "Jan-03", sunrise: "7:16 a.m.", sunset: "4:31 p.m."},
		Day{day: "Jan-05", sunrise: "7:16 a.m.", sunset: "4:33 p.m."},
	}

	var datetests = []struct {
		rawDate time.Time
		date    string
		time    string
		sunrise string
		sunset  string
	}{
		{time.Date(2016, 1, 1, 16, 29, 1, 1, time.UTC), "Jan-01", "4:29 p.m.", "7:15 a.m.", "4:29 p.m."},
		{time.Date(2016, 1, 3, 3, 30, 1, 1, time.UTC), "Jan-03", "3:30 a.m.", "7:16 a.m.", "4:31 p.m."},
		{time.Date(2016, 1, 5, 0, 0, 1, 1, time.UTC), "Jan-05", "12:00 a.m.", "7:16 a.m.", "4:33 p.m."},
	}

	for _, today := range datetests {
		actual := NewToday(today.rawDate, days)
		expected := Today{date: today.date, time: today.time, sunrise: today.sunrise, sunset: today.sunset}
		if actual != expected {
			t.Error(fmt.Sprintf("Expected %v but was %v", expected, actual))
		}
	}
}
