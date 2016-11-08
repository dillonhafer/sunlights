package main

import (
	"errors"
	"fmt"
)

type Day struct {
	day     string
	sunrise string
	sunset  string
}

func findDay(days []Day, today string) (Day, error) {
	for _, day := range days {
		if today == day.day {
			return day, nil
		}
	}

	return Day{}, errors.New(fmt.Sprintf("Could not find entry for '%s' in config.json", today))
}
