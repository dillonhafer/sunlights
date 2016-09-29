package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/savaki/go.hue"
	"log"
	"os"
	"strings"
	"time"
)

const Version = "0.0.1"

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

func findDay(sunsetTable, today string) (Day, error) {
	f, err := os.Open(sunsetTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open sunset file:  %s\n", err)
	}
	r := csv.NewReader(bufio.NewReader(f))
	result, err := r.ReadAll()

	for i := range result {
		day := Day{day: result[i][0], sunrise: result[i][1], sunset: result[i][2]}
		if today == day.day {
			return day, err
		}
	}
	return Day{day: "", sunrise: "", sunset: ""}, err
}

func lightsOff(bridgeAddress, username string) {
	println("Turning lights off")
	bridge := hue.NewBridge(bridgeAddress, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find lights:  %s\n", err)
	}

	for _, light := range lights {
		fmt.Printf("Turned off device => %+v\n", light)
		light.Off()
	}
}

func lightsOn(bridgeAddress, username string) {
	println("Turning lights on")
	bridge := hue.NewBridge(bridgeAddress, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find lights:  %s\n", err)
	}

	for _, light := range lights {
		fmt.Printf("Turned on device => %+v\n", light)
		light.On()
	}
}

func currentTime() string {
	ct := time.Now().Format("3:04 pm")
	ct = strings.Replace(ct, "am", "a.m.", -1)
	ct = strings.Replace(ct, "pm", "p.m.", -1)
	return ct
}

func CurrentDate(t time.Time) string {
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

	date := CurrentDate(time.Now())
	time := currentTime()

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
