package main

import (
	"bufio"
	"encoding/csv"
	"github.com/savaki/go.hue"
	"os"
	"strings"
	"time"
)

var BridgeAddress = os.Getenv("BRIDGE_ADDRESS")
var Username = os.Getenv("USERNAME")
var SunsetTable = os.Getenv("SUNSET_TABLE")

type Day struct {
	day     string
	sunrise string
	sunset  string
}

func findDay() Day {
	f, _ := os.Open(SunsetTable)
	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()

	today := currentDate()
	for i := range result {
		day := Day{day: result[i][0], sunrise: result[i][1], sunset: result[i][2]}
		if today == day.day {
			return day
		}
	}
	return Day{day: "", sunrise: "", sunset: ""}
}

func lightsOff() {
	println("Turning lights off")
	bridge := hue.NewBridge(BridgeAddress, Username)
	lights, _ := bridge.GetAllLights()

	for _, light := range lights {
		light.Off()
	}
}

func lightsOn() {
	println("Turning lights on")
	bridge := hue.NewBridge(BridgeAddress, Username)
	lights, _ := bridge.GetAllLights()

	for _, light := range lights {
		light.On()
	}
}

func currentTime() string {
	ct := time.Now().Format("3:04 pm")
	ct = strings.Replace(ct, "am", "a.m.", -1)
	ct = strings.Replace(ct, "pm", "p.m.", -1)
	return ct
}

func currentDate() string {
	cd := time.Now().Format("Jan-2")
	return cd
}

func main() {
	if BridgeAddress == "" {
		println("You did not set BRIDGE_ADDRESS env var")
		os.Exit(1)
	}

	if Username == "" {
		println("You did not set USERNAME env var")
		os.Exit(1)
	}

	if SunsetTable == "" {
		println("You did not set SUNSET_TABLE env var")
		os.Exit(1)
	}

	day := findDay()
	time := currentTime()

	switch time {
	case day.sunrise:
		lightsOff()
	case day.sunset:
		lightsOn()
	}
}
