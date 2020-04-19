package main

// OnAir: a golang application that will change the color of a Philips Hue light when
// your computer's camera or microphone is enabled. You should take a look at the
// README.md for more information on how to use this application.
//
// watcher.go: This is where you'll find the bits that actually tail the log file and
//             track state. You could have this parse a different log format if you
//             wanted, maybe output something different when video is active vs. inactive.

import (
	"fmt"
	"io"
	"os"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli/v2"
	"github.com/vjeantet/grok"
)

type deviceState struct {
	Name  string
	Type  string
	State string
}

func run(c *cli.Context) {
	var state map[string]deviceState
	state = make(map[string]deviceState)
	config := tail.Config{
		Follow: true,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
	}

	// If Hue UID isnt set, let's assume it's never been set, and do the thing where we help
	// the user out.
	if c.String("hueuid") == "" {
		fmt.Printf("Hue username is not set. Make sure you run onair init or specify the uid on the command line.\n")
		os.Exit(1)
	}

	// tail the log
	t, err := tail.TailFile(c.String("mslog"), config)
	if err != nil {
		fmt.Printf("Cant tail the file: %s\n", err)
		os.Exit(1)
	}

	// configure parser
	gr, _ := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	gr.AddPattern("AMPM", "[AP][M]")
	snitchFormat := "^%{MONTH:month} %{MONTHDAY:day}, %{YEAR:year} at %{HOUR:hour}:%{MINUTE:minute}:%{SECOND:second} %{AMPM:ampm}: %{WORD:devicetype} Device became %{WORD:onoff}: %{GREEDYDATA:device}$"
	gr.AddPattern("SNITCHLOG", snitchFormat)

	// Get connected to the Hue bridge
	bridge, err := loginHue(c.String("hueuid"), c.String("hueip"))
	if err != nil {
		fmt.Printf("Could not log into the Hue bridge: %s", err)
	}

	// Get our xy vals ready
	activex, activey, err := parseXYval(c.String("active"))
	if err != nil {
		fmt.Printf("Could not parse active xy value: %s\n", err)
	}

	inactivex, inactivey, err := parseXYval(c.String("inactive"))
	if err != nil {
		fmt.Printf("Could not parse inactive xy value: %s\n", err)
	}

	// Parse each line of the file. Update state when you do, then check the state to see if anything is active.
	for line := range t.Lines {
		msg, _ := gr.Parse("%{SNITCHLOG}", line.Text)
		if len(msg) == 0 {
			// don't care about devices connected/disconnected
			continue
		}

		state[msg["device"]] = deviceState{
			Type:  msg["deviceType"],
			State: msg["onoff"],
		}

		// read the state of everything. see if anythings active?
		var on bool
		for _, v := range state {
			if v.State == "active" {
				on = true
			}
		}

		// do something with that info
		if on {
			setLights(bridge, c.Int("light"), activex, activey, c.Int("brightness"))
		} else {
			setLights(bridge, c.Int("light"), inactivex, inactivey, c.Int("brightness"))
		}
	}
}
