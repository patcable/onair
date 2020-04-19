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

type lightConfig struct {
	System     string
	Parameters interface{}
}

// run: where the magic happens! this function tails the log file, and for each update to the file, will
// update the system state and then check that state. if any devices are "active" then it will do a thing.
func run(c *cli.Context) {
	var state map[string]deviceState
	state = make(map[string]deviceState)
	config := tail.Config{
		Follow: true,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
	}

	// Do some checks here to make sure we have the bits we need. If you add a lighting system. probably worth adding
	// something like this here. Don't use the `required` flag of urfave/cli since people may be using different light
	// systems.
	if c.String("system") == "hue" && c.String("hueuid") == "" {
		fmt.Printf("Hue username is not set. Make sure you run onair init or specify the uid on the command line.\n")
		os.Exit(1)
	}

	// Get our log file ready for use here
	t, err := tail.TailFile(c.String("log"), config)
	if err != nil {
		fmt.Printf("Cant tail the file: %s\n", err)
		os.Exit(1)
	}

	// Configure parser - you'd add a new log format over in that function.
	// You'll need to do Stuff down in the for loop below to support the different
	// format as well.
	gr, parser, err := configureParser(c.String("logtype"))
	if err != nil {
		fmt.Printf("Cant set up the parser: %s\n", err)
		os.Exit(1)
	}

	// Configure the light - add new lighting systems over there.
	light, err := configureLightSystem(c)
	if err != nil {
		fmt.Printf("Error parsing light config: %s\n", err.Error())
		os.Exit(1)
	}

	// Parse each line of the file. Update state when you do, then check the state to see if anything is active.
	for line := range t.Lines {
		msg, _ := gr.Parse(parser, line.Text)
		if len(msg) == 0 {
			// don't care about devices connected/disconnected
			continue
		}

		state[msg["device"]] = deviceState{
			Type:  msg["deviceType"],
			State: msg["onoff"],
		}

		// read the state of everything. see if anythings active?
		var computerListening bool
		for _, v := range state {
			if v.State == "active" {
				computerListening = true
			}
		}
		err := setLight(light, computerListening)
		if err != nil {
			fmt.Printf("run: Unable to setLightWithContext: %s\n", err)
		}
	}
}

// Set up the grok parser. New formats go here. Make sure to edit the loop under run() too.
func configureParser(logtype string) (gr *grok.Grok, parser string, err error) {
	gr, err = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		return nil, "", err
	}
	gr.AddPattern("AMPM", "[AP][M]")
	snitchFormat := "^%{MONTH:month} %{MONTHDAY:day}, %{YEAR:year} at %{HOUR:hour}:%{MINUTE:minute}:%{SECOND:second} %{AMPM:ampm}: %{WORD:devicetype} Device became %{WORD:onoff}: %{GREEDYDATA:device}$"
	gr.AddPattern("SNITCHLOG", snitchFormat)

	switch logtype {
	case "microsnitch":
		parser = "%{SNITCHLOG}"
	default:
		return nil, "", fmt.Errorf("Your specified a parser of %s which is invalid", logtype)
	}
	return gr, parser, nil
}

// configureLightSystem will... take the config vars and set up the lighting sytem. If you want
// to add a lighting system, this is one of the places you'd need edit to do that.
// Take a look at the hueConfig struct in hue.go - you'll want one of those for your lighting
// system. You'll want to also set up the convenience commands too - ie. `hue init` `hue lights`
// so that you can get the values you'd need for that struct. Finally, update setLightWithContext
// below so that depending on system type, you send the right function for controlling the light.
func configureLightSystem(c *cli.Context) (light lightConfig, err error) {
	switch c.String("system") {
	case "hue":
		// Get connected to the Hue bridge
		bridge, err := loginHue(c.String("hueuid"), c.String("hueip"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("Could not log into the Hue bridge: %s", err)
		}

		// Get our xy vals ready
		activex, activey, err := parseXYval(c.String("hueactive"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("Could not parse active xy value: %s", err)
		}

		inactivex, inactivey, err := parseXYval(c.String("hueinactive"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("Could not parse inactive xy value: %s", err)
		}

		light = lightConfig{
			System: "hue",
			Parameters: hueConfig{
				Bridge:     bridge,
				Light:      c.Int("huelight"),
				Brightness: c.Int("huebrightness"),
				Active:     []float32{activex, activey},
				Inactive:   []float32{inactivex, inactivey},
			},
		}
	default:
		return lightConfig{}, fmt.Errorf("You specified a lighting system of %s which is invalid", c.String("system"))
	}
	return light, nil
}

// setLightWithContext is where you can actually send the command to do a thing based on the light system you use.
// lightConfig.Parameters is an interface, so typecast that then you can shoot off whatever function you want.
func setLight(light lightConfig, computerListening bool) (err error) {
	switch light.System {
	case "hue":
		settings := light.Parameters.(hueConfig)
		if computerListening {
			err = setHueLights(settings.Bridge, settings.Light, settings.Active[0], settings.Active[1], settings.Brightness)
		} else {
			err = setHueLights(settings.Bridge, settings.Light, settings.Inactive[0], settings.Inactive[1], settings.Brightness)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
