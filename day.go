package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type Day struct {
	day     string
	sunrise string
	sunset  string
}

func findDay(sunsetTable, today string) (Day, error) {
	csvFile, err := os.Open(sunsetTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open sunset file: %s\n", err)
		return Day{}, err
	}

	csvRows := csv.NewReader(bufio.NewReader(csvFile))
	result, err := csvRows.ReadAll()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read csv file: %s\n", err)
		return Day{}, err
	}

	for i := range result {
		row := result[i]
		day := Day{day: row[0], sunrise: row[1], sunset: row[2]}

		if today == day.day {
			return day, nil
		}
	}

	return Day{}, errors.New(fmt.Sprintf("Could not find entry for '%s' in csv", today))
}
