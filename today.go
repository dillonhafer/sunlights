package main

import (
	"log"
	"strings"
	"time"
)

type Today struct {
	date    string
	time    string
	sunrise string
	sunset  string
}

func NewToday(currentTime time.Time, days []Day) Today {
	date := FormatDate(currentTime)
	time := FormatTime(currentTime)
	day, err := findDay(days, date)
	if err != nil {
		log.Fatal(err)
	}

	return Today{date: date, time: time, sunrise: day.Sunrise, sunset: day.Sunset}
}

func FormatTime(t time.Time) string {
	ct := t.Format("3:04 pm")
	ct = strings.Replace(ct, "am", "a.m.", -1)
	ct = strings.Replace(ct, "pm", "p.m.", -1)
	return ct
}

func FormatDate(t time.Time) string {
	return t.Format("Jan-02")
}
