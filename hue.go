package main

// OnAir: a golang application that will change the color of a Philips Hue light when
// your computer's camera or microphone is enabled. You should take a look at the
// README.md for more information on how to use this application.
//
// hue.go: a collection of functions for actually having the Hue change colors or
//         brightness. If you want to modify what information `onair lights` displays
//         or update how the light updates, you're in the right place.

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amimof/huego"
	"github.com/urfave/cli/v2"
)

type hueConfig struct {
	Bridge     *huego.Bridge
	Light      int
	Brightness int
	Active     []float32
	Inactive   []float32
}

func initHue(c *cli.Context) {
	if c.String("hueuid") != "" {
		fmt.Printf("Seems like you already have a username set. Remove that variable from your config file if that isn't the case.\n")
		os.Exit(1)
	}

	// find the bridge
	bridge, err := huego.Discover()
	if err != nil {
		fmt.Printf("Could not discover the Hue bridge.\n")
	}

	fmt.Printf("Head on over to your HUE Bridge and push the button to create a new user.\n")
	fmt.Printf("Waiting %d seconds - you can set this longer by passing --timeout to 'onair init'.\n", c.Int("timeout"))
	time.Sleep(time.Duration(c.Int("timeout")) * time.Second)
	user, err := bridge.CreateUser("OnAir")
	if err != nil {
		fmt.Printf("Could not create user: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("   Username: %s\n", user)
	fmt.Printf("Save this value, you'll need it for the config file.\n")
}

func loginHue(username string, ip string) (bridge *huego.Bridge, err error) {
	if ip == "" {
		bridge, err = huego.Discover()
		if err != nil {
			return nil, fmt.Errorf("Could not discover the Hue bridge - is https://www.meethue.com/api/nupnp accessible?")
		}

		bridge = bridge.Login(username)
	} else {
		bridge = huego.New(ip, username)
	}
	return bridge, nil
}

func getHueLights(c *cli.Context) {
	bridge, err := loginHue(c.String("hueuid"), c.String("hueip"))
	if err != nil {
		fmt.Printf("Could not login to bridge: %s\n", err)
		os.Exit(1)
	}

	lights, err := bridge.GetLights()
	if err != nil {
		fmt.Printf("Could not get lights: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your Hue bridge knows about these lights:\n")
	for _, v := range lights {
		if v.State.On {
			fmt.Printf(" -> %s:\n", v.Name)
			fmt.Printf("      ID:         %d\n", v.ID)
			fmt.Printf("      Color Mode: %s\n", v.State.ColorMode)
			fmt.Printf("      Brightness: %d\n", v.State.Bri)
			fmt.Printf("      Hue:        %d\n", v.State.Hue)
			fmt.Printf("      Saturation: %d\n", v.State.Sat)
			fmt.Printf("      Color Temp: %d\n", v.State.Ct)
			fmt.Printf("      XY:         %v\n", v.State.Xy)
		} else {
			fmt.Printf(" -> %s (off)\n", v.Name)
		}
	}
}

func setHueLights(bridge *huego.Bridge, light int, x float32, y float32, bri int) error {
	_, err := bridge.GetLight(light)
	if err != nil {
		return err
	}

	newState := huego.State{
		On:  true,
		Bri: uint8(bri),
		Xy:  []float32{x, y},
	}

	_, err = bridge.SetLightState(light, newState)
	if err != nil {
		return err
	}

	return nil
}

func parseXYval(xyval string) (float32, float32, error) {
	split := strings.Split(xyval, ",")
	convx, err := strconv.ParseFloat(split[0], 32)
	if err != nil {
		return 0, 0, err
	}
	convy, err := strconv.ParseFloat(split[1], 32)
	if err != nil {
		return 0, 0, err
	}
	return float32(convx), float32(convy), nil
}
