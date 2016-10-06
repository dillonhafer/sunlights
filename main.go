package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	hue "github.com/dillonhafer/go.hue"
)

const Version = "2.0.0"

var options struct {
	setup         bool
	version       bool
	app           string
	bridgeAddress string
	username      string
	sunsetTable   string
}

type Day struct {
	day     string
	sunrise string
	sunset  string
}

func puts(message string) {
	l := log.New(os.Stdout, "[Sunlights] ", log.Ldate|log.Ltime)
	l.Printf("%s", message)
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

	return Day{}, errors.New("Could not find entry for '%s' in csv")
}

func lightsOff(bridgeAddress, username string) {
	puts("Turning lights off")

	bridge := hue.NewBridge(bridgeAddress, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find lights:  %s\n", err)
	}

	for _, light := range lights {
		puts(fmt.Sprintf("Turned off light => %+v\n", light.Name))
		light.Off()
	}
}

func lightsOn(bridgeAddress, username string) {
	puts("Turning lights on")
	bridge := hue.NewBridge(bridgeAddress, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		puts(fmt.Sprintf("Could not find lights: %s\n", err))
	}

	for _, light := range lights {
		puts(fmt.Sprintf("Turned on light => %+v\n", light.Name))
		light.On()
	}
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

func setup(app string) {
	locators, err := hue.DiscoverBridges(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not discover bridge:  %s\n", err)
	}
	locator := locators[0]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Push the button on your Hue and the press any key to continue...")
	reader.ReadString('\n')

	bridge, err := locator.CreateUser(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not discover bridge:  %s\n", err)
	}

	fmt.Printf("registered new device => %+v\n", bridge)
	os.Exit(0)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.sunsetTable, "sunsetTable", "", "path to CSV of time data (e.g. /times.csv)")
	flag.StringVar(&options.username, "username", "", "username use to connect to bridge")
	flag.StringVar(&options.bridgeAddress, "bridgeAddress", "", "i.p. address of Hue bridge (e.g. 192.168.2.2)")
	flag.StringVar(&options.app, "app", "", "name of app. (Only used if -setup flag is pressent)")
	flag.BoolVar(&options.setup, "setup", false, "set this flag to setup a new Hue bridge")
	flag.BoolVar(&options.version, "version", false, "print version and exit")
	flag.Parse()

	if options.setup {
		if options.app == "" {
			println("You did not set `-app` name to use during setup.")
			os.Exit(1)
		}

		setup(options.app)
	}

	if options.bridgeAddress == "" {
		println("You did not set a bridgeAddress flag")
		os.Exit(1)
	}

	if options.username == "" {
		println("You did not set a username flag")
		os.Exit(1)
	}

	if options.sunsetTable == "" {
		println("You did not set a sunsetTable flag")
		os.Exit(1)
	}

	ct := time.Now()
	date := FormatDate(ct)
	time := FormatTime(ct)

	day, err := findDay(options.sunsetTable, date)
	if err != nil {
		log.Fatal(err)
	}

	switch time {
	case day.sunrise:
		lightsOff(options.bridgeAddress, options.username)
	case day.sunset:
		lightsOn(options.bridgeAddress, options.username)
	}
}
