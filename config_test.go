package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	hue "github.com/dillonhafer/go.hue"
)

func TestWrite(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "test-config-file")
	file := configFile.Name()
	defer os.Remove(file)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	bulbs := []LightBulb{}
	bulbs = append(bulbs, LightBulb{Name: "Bedroom"})

	config := Config{file: file, Username: "sunlightuser1", BridgeAddress: "192.168.1.1", LightBulbs: bulbs}
	config.write()

	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Fatal("Config file not written")
	}

	c, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal("Could not read config file")
	}

	expectedConfig := `{
  "BridgeAddress": "192.168.1.1",
  "Username": "sunlightuser1",
  "LightBulbs": [
    {
      "Name": "Bedroom"
    }
  ],
  "Days": null
}`

	if string(c) != expectedConfig {
		t.Fatalf("File not written properly, got: %v", string(c))
	}
}

func TestLoad(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "test-config-file")
	file := configFile.Name()
	defer os.Remove(file)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	contents := []byte(`{"BridgeAddress":"192.168.1.1","Username":"sunlightuser1","Days":null,"LightBulbs":null}`)
	ioutil.WriteFile(file, contents, 0644)

	config := Config{file: file}
	config.Load()

	if config.Username != "sunlightuser2" && config.BridgeAddress != "192.168.1.1" {
		t.Fatal(fmt.Sprintf("Config was not loaded: %v", config))
	}
}

func TestAddBulb(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "test-config-file")
	file := configFile.Name()
	defer os.Remove(file)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	config := Config{file: file}
	config.AddBulb("Bedroom")

	expectedConfig := `{
  "BridgeAddress": "",
  "Username": "",
  "LightBulbs": [
    {
      "Name": "Bedroom"
    }
  ],
  "Days": null
}`

	c, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal("Could not read config file")
	}

	if string(c) != expectedConfig {
		t.Fatalf("Bulb not added, got: %v", string(c))
	}

	if config.LightBulbs[0].Name != "Bedroom" {
		t.Fatalf("Bulb not in config, got: %v", config.LightBulbs)
	}
}

func TestRemoveBulb(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "test-config-file")
	file := configFile.Name()
	defer os.Remove(file)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	config := Config{file: file, Username: "test", BridgeAddress: "localhost"}
	config.AddBulb("Bedroom")
	config.AddBulb("Living Room")

	if len(config.LightBulbs) != 2 {
		t.Fatalf("Test setup not correct, got: %v", config.LightBulbs)
	}

	config.RemoveBulb("Living Room")

	expectedConfig := `{
  "BridgeAddress": "localhost",
  "Username": "test",
  "LightBulbs": [
    {
      "Name": "Bedroom"
    }
  ],
  "Days": null
}`

	c, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal("Could not read config file")
	}

	if string(c) != expectedConfig {
		t.Fatalf("Bulb not removed, got: %v", string(c))
	}

	if len(config.LightBulbs) == 2 {
		t.Fatalf("Bulb still in config, got: %v", config.LightBulbs)
	}
}

func TestSaveSetup(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "test-config-file")
	file := configFile.Name()
	defer os.Remove(file)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	bridge := &hue.Bridge{Username: "test", IpAddr: "localhost"}
	config := Config{file: file}
	config.SaveSetup(bridge)

	expectedConfig := `{
  "BridgeAddress": "localhost",
  "Username": "test",
  "LightBulbs": null,
  "Days": null
}`

	c, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal("Could not read config file")
	}

	if string(c) != expectedConfig {
		t.Fatalf("Setup not saved, got: %v", string(c))
	}

	if config.Username != "test" {
		t.Fatalf("Username not set to 'test', got: %v", config.Username)
	}

	if config.BridgeAddress != "localhost" {
		t.Fatalf("Username not set to 'localhost', got: %v", config.BridgeAddress)
	}
}
