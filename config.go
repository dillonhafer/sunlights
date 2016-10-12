package main

import (
	"encoding/csv"
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
