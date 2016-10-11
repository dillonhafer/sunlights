package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	hue "github.com/dillonhafer/go.hue"
	_ "github.com/mattn/go-sqlite3"
)

const Version = "2.0.0"

var Bridge *hue.Bridge

type LightBulb struct {
	Id   int
	Name string
}

var options struct {
	setup         bool
	version       bool
	app           string
	bridgeAddress string
	username      string
	sunsetTable   string
}

func puts(message string) {
	l := log.New(os.Stdout, "[Sunlights] ", log.Ldate|log.Ltime)
	l.Printf("%s", message)
}

func toggleLights(on bool) {
	direction := "on"
	if !on {
		direction = "off"
	}
	puts(fmt.Sprintf("Turning lights %s", direction))

	lights, err := Bridge.GetAllLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find lights:  %s\n", err)
	}

	var allowedLights []*hue.Light
	for _, light := range lights {
		for _, l := range AllBulbs() {
			if l.Name == light.Name {
				allowedLights = append(allowedLights, light)
			}
		}
	}

	for _, light := range allowedLights {
		puts(fmt.Sprintf("Turned %s light => %+v\n", direction, light.Name))
		if on {
			light.On()
		} else {
			light.Off()
		}
	}
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

func addBulb(name string) {
	db, err := sql.Open("sqlite3", "./sunlights.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into light_bulbs(name) values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
	println(fmt.Sprintf("Added %s", name))
}

func AllBulbs() []LightBulb {
	db, err := sql.Open("sqlite3", "./sunlights.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select id, name from light_bulbs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bulbs []LightBulb
	for rows.Next() {
		lightBulb := LightBulb{}
		err = rows.Scan(&lightBulb.Id, &lightBulb.Name)
		if err != nil {
			log.Fatal(err)
		}
		bulbs = append(bulbs, lightBulb)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return bulbs
}

func listBulbs() {
	println("\n", "Here are your light bulbs:")
	bulbs := AllBulbs()
	for _, light := range bulbs {
		fmt.Println(fmt.Sprintf("   %d.", light.Id), light.Name)
	}
	println()
}

func setupDatabase() {
	db, err := sql.Open("sqlite3", "./sunlights.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	  create table light_bulbs (id integer not null primary key, name text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func lookForSubCommands(args []string) {
	if len(args) > 1 {
		switch args[1] {
		case "list", "ls":
			listBulbs()
			os.Exit(0)
		case "add", "a":
			if len(args) > 2 {
				addBulb(args[2])
			} else {
				println("You must provide a name for your new light bulb")
			}
			os.Exit(0)
		case "db:setup":
			setupDatabase()
			println("Created database")
			os.Exit(0)
		}
	}
}

func main() {
	lookForSubCommands(os.Args)

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

	Bridge = hue.NewBridge(options.bridgeAddress, options.username)
	today := NewToday(time.Now(), options.sunsetTable)

	switch today.time {
	case today.sunrise:
		off := false
		toggleLights(off)
	case today.sunset:
		on := true
		toggleLights(on)
	}
}
