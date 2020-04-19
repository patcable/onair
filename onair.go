package main

// OnAir: a golang application that will change the color of a Philips Hue light when
// your computer's camera or microphone is enabled. You should take a look at the
// README.md for more information on how to use this application.
//
// onair.go: cli skeleton and flag definitions. You can find some default values here,
//           so if you'd like to change those here's a likely place to make that happen.

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Printf("onair requires $HOME to be set")
		os.Exit(1)
	}
	// Some defaults.?
	microSnitchLogFile := fmt.Sprintf("%s/Library/Logs/Micro Snitch.log", homeDir)
	configFile := fmt.Sprintf("%s/.onair.yml", homeDir)

	globalFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "mslog",
			Aliases: []string{"m"},
			Usage:   "Path to the Micro Snitch Log File",
			Value:   microSnitchLogFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "hueuid",
			Aliases: []string{"u"},
			Usage:   "The username to use for the Hue bridge. Pulls from $HOME/.onair by default",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "hueip",
			Aliases: []string{"i"},
			Usage:   "The IP of your Hue bridge - onair uses automatic discovery if not specified",
		}),
		&cli.StringFlag{
			Name:  "config",
			Usage: "Path to the config yaml file",
			Value: configFile,
		},
	}

	runFlags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "light",
			Aliases: []string{"l"},
			Usage:   "ID of the light you'll control",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "brightness",
			Aliases: []string{"b"},
			Usage:   "Brightness (1-254)",
			Value:   70,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "active",
			Usage: "XY color value when video/audio is active",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "inactive",
			Usage: "XY color value when video/audio is inactive",
		}),
		&cli.StringFlag{
			Name:   "config",
			Usage:  "Path to the config yaml file",
			Value:  configFile,
			Hidden: true,
		},
	}

	app := &cli.App{
		Name:   "onair",
		Usage:  "monitors your audio/video devices and controls a Philips HUE light based on their status",
		Before: altsrc.InitInputSourceWithContext(globalFlags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:  globalFlags,
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Set up onair for the first time",
				Action: func(c *cli.Context) error {
					initHue(c)
					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "timeout",
						Usage: "how long to wait for you to push the button on your Hue bridge",
						Value: 15,
					},
				},
			},
			{
				Name:  "lights",
				Usage: "display information about available lights and their color settings",
				Action: func(c *cli.Context) error {
					getHueLights(c)
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "run the log watcher/set your lights",
				Action: func(c *cli.Context) error {
					run(c)
					return nil
				},
				Before: altsrc.InitInputSourceWithContext(runFlags, altsrc.NewYamlSourceFromFlagFunc("config")),
				Flags:  runFlags,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("oh no: %s\n", err)
	}
}
