package main

// OnAir: a golang application that will change the color of a Philips Hue light when
// your computer's camera or microphone is enabled. You should take a look at the
// README.md for more information on how to use this application.
//
// iftt.go: function to invoke IFTTT hooks
//      

import (
	"fmt"
	"net/http"
)

type iftttConfig struct {
	key        string
	onairHook  string
	offairHook  string
}

func invokeIFTTTHook(key string, hook string) error {
       fmt.Printf("invokes key %s hook %s\n", key, hook)
       url := "https://maker.ifttt.com/trigger/" + hook + "/with/key/" + key
       fmt.Printf("url: %s \n", url)
       resp, err := http.Get(url)
       _ = resp   // ignore the response
       return err
}

