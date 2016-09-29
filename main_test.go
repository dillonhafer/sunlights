package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func assertEqual(t *testing.T, given Day, expected Day) {
	if given != expected {
		t.Fatalf("\033[31mExpected \033[m \033[33m%v\033[33m \033[31mbut was\033[m \033[33m%v\033[m", expected, given)
	}
}

func createTestFile() string {
	content := []byte(`day,sunrise,sunset
Sep-21,6:33 a.m.,6:45 p.m.
`)

	tmpfile, err := ioutil.TempFile("", "tmp-sunset.csv")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func TestFindDay(t *testing.T) {
	expectedDay := Day{day: "Sep-21", sunrise: "6:33 a.m.", sunset: "6:45 p.m."}
	testFile := createTestFile()
	defer os.Remove(testFile)

	day := findDay(testFile, "Sep-21")

	assertEqual(t, day, expectedDay)
}
