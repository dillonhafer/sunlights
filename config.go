package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	hue "github.com/dillonhafer/go.hue"
)

type LightBulb struct {
	Name string
}

type Config struct {
	file          string
	BridgeAddress string
	Username      string
	LightBulbs    []LightBulb
	Days          []Day
}

func (c *Config) write() {
	file, err := os.Create(c.file)
	if err != nil {
		log.Fatal("Cannot create config file", err)
	}
	defer file.Close()

	configJson, _ := json.MarshalIndent(c, "", "  ")
	err = ioutil.WriteFile(c.file, configJson, 0644)
	if err != nil {
		log.Fatal("Cannot write to config file", err)
	}
}

func (c *Config) Load() {
	file, err := ioutil.ReadFile(c.file)
	if err != nil {
		log.Fatal(errors.New("Config file does not exist. Please run `sunlights setup`"))
	}

	if err := json.Unmarshal(file, c); err != nil {
		log.Fatal(err)
	}

	if c.Username != "" && c.BridgeAddress != "" {
		return
	}

	log.Fatal(errors.New("Config file is missing data. Please run `sunlights setup` again"))
}

func (c *Config) AddBulb(name string) {
	c.LightBulbs = append(c.LightBulbs, LightBulb{Name: name})
	c.write()
}

func (c *Config) RemoveBulb(name string) {
	var keepBulbs []LightBulb
	for _, bulb := range c.LightBulbs {
		if bulb.Name != name {
			keepBulbs = append(keepBulbs, bulb)
		}
	}
	c.LightBulbs = keepBulbs
	c.write()
}

func (c *Config) SaveSetup(b *hue.Bridge) {
	c.Username = b.Username
	c.BridgeAddress = b.IpAddr
	c.write()
}
