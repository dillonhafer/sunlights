package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"log"
	"os"
)

type Config struct {
	BridgeAddress string
	Username      string
}

func (c *Config) Write(fileName string) {
	if fileName == "" {
		fileName = "./config.csv"
	}
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create config file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	data := [][]string{{"bridgeAddress", "username"}, {c.BridgeAddress, c.Username}}
	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write to config file", err)
		}
	}

	defer writer.Flush()
}

func (c *Config) Fetch(fileName string) error {
	if fileName == "" {
		fileName = "./config.csv"
	}
	csvFile, err := os.Open(fileName)
	if err != nil {
		return errors.New("Config file does not exist. Please run `sunlights setup`")
	}

	csvRows := csv.NewReader(bufio.NewReader(csvFile))
	result, err := csvRows.ReadAll()

	if err != nil {
		return err
	}

	for i := range result {
		row := result[i]
		c.BridgeAddress = row[0]
		c.Username = row[1]
	}

	if c.Username != "" && c.BridgeAddress != "" {
		return nil
	}

	return errors.New("Config file is missing data. Please run `sunlights setup` again")
}
