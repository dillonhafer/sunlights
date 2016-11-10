package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	hue "github.com/dillonhafer/go.hue"
	"github.com/urfave/cli"
)

const Version = "3.1.0"

var config = Config{}

func puts(message string) {
	l := log.New(os.Stdout, "[Sunlights] ", log.Ldate|log.Ltime)
	l.Printf("%s", message)
}

func setup() {
	locators, err := hue.DiscoverBridges(false)
	if err != nil || len(locators) == 0 {
		log.Fatal("Could not discover bridge. Make sure you are on the same network as the bridge.")
	}
	locator := locators[0]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Push the button on your Hue and the press any key to continue...")
	reader.ReadString('\n')

	bridge, err := locator.CreateUser("Sunlights")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not discover bridge:  %s\n", err)
	}

	config.SaveSetup(bridge)
	fmt.Printf("Setup is complete")
	os.Exit(0)
}

func AllBulbs() []LightBulb {
	if len(config.LightBulbs) == 0 {
		println("It looks like you haven't added any bulbs yet. Run `sunlights -h` for instructions.")
		os.Exit(1)
	}

	return config.LightBulbs
}

func ActionWithConfig(f func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		config.Load()
		return f(c)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "sunlights"
	app.Usage = "Control lights based on sunrise/sunset"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,C",
			Value:       ".sunlights.json",
			Usage:       "Configuration file",
			Destination: &config.file,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "Pair with bridge for the first time.",
			Action: func(c *cli.Context) error {
				setup()
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all light bulbs.",
			Action: ActionWithConfig(
				func(c *cli.Context) error {
					bulbs := AllBulbs()
					println("\n", " Here are your light bulbs💡")
					for i, bulb := range bulbs {
						fmt.Println(fmt.Sprintf("    %d. %s", i+1, bulb.Name))
					}
					println()
					os.Exit(0)
					return nil
				},
			),
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a light bulb to be controlled",
			Action: ActionWithConfig(
				func(c *cli.Context) error {
					name := c.Args().First()
					if name != "" {
						config.AddBulb(name)
						println(fmt.Sprintf("Added %s", name))
						os.Exit(0)
						return nil
					} else {
						return errors.New("You must provide a name for your new light bulb")
					}
				},
			),
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Remove control of a light bulb",
			Action: ActionWithConfig(
				func(c *cli.Context) error {
					name := c.Args().First()
					if name != "" {
						config.RemoveBulb(name)
						os.Exit(0)
						return nil
					} else {
						return errors.New("You must provide a name for the light bulb you want to remove.")
					}
				},
			),
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show all lightbulbs connected to bridge",
			Action: ActionWithConfig(
				func(c *cli.Context) error {
					bridge := Bridge{Hue: hue.NewBridge(config.BridgeAddress, config.Username)}
					lights, _ := bridge.Hue.GetAllLights()

					for _, light := range lights {
						println(light.Name)
					}
					return nil
				},
			),
		},
	}

	app.Action = ActionWithConfig(
		func(c *cli.Context) error {
			bridge := Bridge{Hue: hue.NewBridge(config.BridgeAddress, config.Username)}
			today := NewToday(time.Now(), config.Days)

			switch today.time {
			case today.sunrise:
				bridge.TurnLightsOff()
			case today.sunset:
				bridge.TurnLightsOn()
			}
			return nil
		},
	)

	app.Run(os.Args)
}
