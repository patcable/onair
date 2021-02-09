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
	"net/http"
	"github.com/urfave/cli/v2"
)

type iftttConfig struct {
	key        string
	onairHook  string
	offairHook  string
}

func checkIFTTTHooks(c *cli.Context) {
  // fill in
}

func invokeIFTTTHook(key string, hook string) error {
       fmt.Printf("invokes key %s hook %s\n", key, hook)
       url := "https://maker.ifttt.com/trigger/" + hook + "/with/key/" + key
       fmt.Printf("url: %s \n", url)
       resp, err := http.Get(url)
       _ = resp   // ignore the response
       return err
}

