package main

import (
	"fmt"

	hue "github.com/dillonhafer/go.hue"
)

type Bridge struct {
	Hue *hue.Bridge
}

func (b *Bridge) TurnLightsOn() {
	on := true
	b.toggleLights(on)
}

func (b *Bridge) TurnLightsOff() {
	off := false
	b.toggleLights(off)
}

func (b *Bridge) toggleLights(on bool) {
	var direction string
	if on {
		direction = "on"
	} else {
		direction = "off"
	}

	puts(fmt.Sprintf("Turning lights %s", direction))

	var allowedLights []*hue.Light
	for _, l := range AllBulbs() {
		light, err := b.Hue.FindLightByName(l.Name)
		if err != nil {
			puts(fmt.Sprintf("Could not find light: %v", err))
		} else {
			allowedLights = append(allowedLights, light)
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
