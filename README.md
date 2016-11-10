Sunlights [![CircleCI](https://circleci.com/gh/dillonhafer/sunlights.svg?style=svg)](https://circleci.com/gh/dillonhafer/sunlights)
------

A cross-platform server for Raspberry Pi, Linux, OSX, or Windows to turn Hue lights on and off based on the sunrise/sunset. Ideally used for outdoor porch lighting.

Hardware I used for project:

1. [Phillips Hue Starter Kit (2nd gen)](http://amzn.to/2cRQwVt)

## Quick setup

1. Run `sunlights setup` (You will need to press the button on the Phillips HUE bridge for this step)
1. See available lightbulbs `sunlights show`
1. Add a light bulb to control `sunlights add 'Hue white lamp 1'`
1. Add time schedules to the `Days` array in `.sunlights.json` config file

```json
{
  "BridgeAddress": "192.168.1.16",
  "Username": "1234",
  "LightBulbs": [
    {"Name": "Hue white lamp 1"}
  ],
  "Days": [
    {
      "Date": "Jan-01",
      "Sunrise": "7:15 a.m.",
      "Sunset": "4:29 p.m."
    }
  ]
}
```

Once the above setup is done, simply add a crontab to check the lights every minute:

```shell
# crontab
* * * * * /path/to/sunlights -C /path/to/.sunlights.json
```

## About

```
NAME:
   sunlights

USAGE:
   sunlights [global options] command [command options] [arguments...]
   
VERSION:
   3.1.0
   
COMMANDS:
     setup       Pair with bridge for the first time.
     list, ls    List all light bulbs.
     add, a      Add a light bulb to be controlled
     remove, rm  Remove control of a light bulb
     show, s     Show all lightbulbs connected to bridge
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -C value  Configuration file (default: ".sunlights.json")
   --help, -h                show help
   --version, -v             print the version
```

## License

Copyright © 2016 Dillon Hafer. It is free software, and may be redistributed under the terms specified below.

**This work is dedicated to the public domain. See UNLICENSE**

**This software uses 3rd-party software - see 3rd-party-licenses for their respective licenses**

