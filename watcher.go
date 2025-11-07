package main

// OnAir: a golang application that will change the color of a Philips Hue light when
// your computer's camera or microphone is enabled. You should take a look at the
// README.md for more information on how to use this application.
//
// watcher.go: This is where you'll find the bits that actually tail the log file and
//             track state. You could have this parse a different log format if you
//             wanted, maybe output something different when video is active vs. inactive.

import (
	"context"
	"fmt"
	"os"
	"time"

	mediaDevices "github.com/patcable/go-media-devices-state"
	"github.com/urfave/cli/v3"
)

type lightConfig struct {
	System     string
	Parameters interface{}
}

// run: where the magic happens! this function tails the log file, and for each update to the file, will
// update the system state and then check that state. if any devices are "active" then it will do a thing.
func run(ctx context.Context, c *cli.Command) {
	// Do some checks here to make sure we have the bits we need. If you add a lighting system. probably worth adding
	// something like this here. Don't use the `required` flag of urfave/cli since people may be using different light
	// systems.
	if c.String("system") == "hue" && c.String("hueuid") == "" {
		fmt.Printf("Hue username is not set. Make sure you run onair init or specify the uid on the command line.\n")
		os.Exit(1)
	}

	// Configure the light - add new lighting systems over there.
	light, err := configureLightSystem(ctx, c)
	if err != nil {
		fmt.Printf("Error parsing light config: %s\n", err.Error())
		os.Exit(1)
	}

	// every n seconds, check the status of the system camera and microphone. If either are hot, change the light.
	for {
		cam, err := mediaDevices.IsCameraOn(c.Bool("debug"))
		if err != nil {
			fmt.Printf("Error with IsCameraOn: %s\n", err.Error())
			os.Exit(1)
		}
		mic, err := mediaDevices.IsMicrophoneOn(c.Bool("debug"))
		if err != nil {
			fmt.Printf("Error with IsMicrophoneOn: %s\n", err.Error())
			os.Exit(1)
		}

		if cam || mic {
			err := setLight(light, true)
			if err != nil {
				fmt.Printf("run: Unable to setLightWithContext: %s\n", err)
			}
		} else {
			err := setLight(light, false)
			if err != nil {
				fmt.Printf("run: Unable to setLightWithContext: %s\n", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// configureLightSystem will take the config vars and set up the lighting system. If you want
// to add a lighting system, this is one of the places you'd need edit to do that.
// Take a look at the hueConfig struct in hue.go - you'll want one of those for your lighting
// system. You'll want to also set up the convenience commands too - ie. `hue init` `hue lights`
// so that you can get the values you'd need for that struct. Finally, update setLightWithContext
// below so that depending on system type, you send the right function for controlling the light.
func configureLightSystem(ctx context.Context, c *cli.Command) (light lightConfig, err error) {
	switch c.String("system") {
	case "hue":
		// Get connected to the Hue bridge
		bridge, err := loginHue(c.String("hueuid"), c.String("hueip"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("could not log into the Hue bridge: %s", err)
		}

		// Get our xy vals ready
		activex, activey, err := parseXYval(c.String("hueactive"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("could not parse active xy value: %s", err)
		}

		inactivex, inactivey, err := parseXYval(c.String("hueinactive"))
		if err != nil {
			return lightConfig{}, fmt.Errorf("could not parse inactive xy value: %s", err)
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
	case "ifttt":
		light = lightConfig{
			System: "ifttt",
			Parameters: iftttConfig{
				key:        c.String("ifttt-key"),
				onairHook:  c.String("ifttt-onair"),
				offairHook: c.String("ifttt-offair"),
			},
		}

	default:
		return lightConfig{}, fmt.Errorf("you specified a lighting system of %s which is invalid", c.String("system"))
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
			if err != nil {
				fmt.Printf("could set hue lights: %s", err)
			}
		} else {
			err = setHueLights(settings.Bridge, settings.Light, settings.Inactive[0], settings.Inactive[1], settings.Brightness)
			if err != nil {
				fmt.Printf("could set hue lights: %s", err)
			}
		}
	case "ifttt":
		settings := light.Parameters.(iftttConfig)
		if computerListening {
			err = invokeIFTTTHook(settings.key, settings.onairHook)
		} else {
			err = invokeIFTTTHook(settings.key, settings.offairHook)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
