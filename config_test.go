package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	configFile := "./test-config.csv"
	defer os.Remove(configFile)

	if _, err := os.Stat(configFile); err == nil {
		err := os.Remove(configFile)

		if err != nil {
			fmt.Println(err)
			return
		}

	}

	config := Config{Username: "sunlightuser1", BridgeAddress: "192.168.1.1"}
	config.Write(configFile)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file not written")
	}
	c, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatal("Could not read config file")
	}

	expectedConfig := `bridgeAddress,username
192.168.1.1,sunlightuser1
`

	if string(c) != expectedConfig {
		t.Fatal("File not written properly")
	}
}

func TestFetch(t *testing.T) {
	configFile := "./test-config.csv"
	file := []byte(`bridgeAddress,username
192.168.1.1,sunlightuser2
`)
	ioutil.WriteFile(configFile, file, 0644)
	defer os.Remove(configFile)

	config := Config{}
	config.Fetch(configFile)

	if config.Username != "sunlightuser2" && config.BridgeAddress != "192.168.1.1" {
		t.Fatal(fmt.Sprintf("Config was not loaded: %v", config))
	}
}
