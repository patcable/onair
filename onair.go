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

	"context"

	altsrc "github.com/urfave/cli-altsrc/v3"
	altsrcyaml "github.com/urfave/cli-altsrc/v3/yaml"
	"github.com/urfave/cli/v3"
)

var buildVersion string
var configFile string

func main() {
	if os.Getenv("ONAIR_CONFIG") == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			fmt.Printf("onair requires $HOME to be set")
			os.Exit(1)
		}
		configFile = fmt.Sprintf("%s/.onair.yml", homeDir)
	} else {
		configFile = os.Getenv("ONAIR_CONFIG")
	}

	globalFlags := []cli.Flag{
		&cli.IntFlag{
			Name:    "interval",
			Aliases: []string{"t"},
			Usage:   "interval in seconds to poll the CoreMediaIO system",
			Value:   2,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "print detected cameras/mics on every poll",
			Value:   false,
		},
	}

	hueFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "hueuid",
			Aliases: []string{"u"},
			Usage:   "username for the Hue bridge",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueuid", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:    "hueip",
			Aliases: []string{"i"},
			Usage:   "ip address of the Hue bridge (automatic discovery if not specified)",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueip", altsrc.StringSourcer(configFile)),
			),
		},
	}

	runFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "system",
			Aliases: []string{"s"},
			Usage:   "which light system do you use (currently only supports hue and ifttt webhooks)",
			Value:   "hue",
		},
		&cli.StringFlag{
			Name:    "hueuid",
			Aliases: []string{"u"},
			Usage:   "username for the Hue bridge",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueuid", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:    "hueip",
			Aliases: []string{"i"},
			Usage:   "ip address of the Hue bridge (automatic discovery if not specified)",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueip", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.IntFlag{
			Name:    "huelight",
			Aliases: []string{"l"},
			Usage:   "ID of the light you'll control",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("huelight", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.IntFlag{
			Name:    "huebrightness",
			Aliases: []string{"b"},
			Usage:   "Brightness (1-254)",
			Value:   70,
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("huebrightness", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:  "hueactive",
			Usage: "XY color value when video/audio is active",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueactive", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:  "hueinactive",
			Usage: "XY color value when video/audio is inactive",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("hueinactive", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:    "ifttt-key",
			Aliases: []string{"k"},
			Usage:   "key for IFTTT webhook requests",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("ifttt-key", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:    "ifttt-onair",
			Aliases: []string{"o"},
			Usage:   "Name of IFTTT webhook invoked when video/audio becomes active",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("ifttt-onair", altsrc.StringSourcer(configFile)),
			),
		},
		&cli.StringFlag{
			Name:    "ifttt-offair",
			Aliases: []string{"f"},
			Usage:   "Name of IFTTT webhook invoked when video/audio becomes inactive",
			Sources: cli.NewValueSourceChain(
				altsrcyaml.YAML("ifttt-offair", altsrc.StringSourcer(configFile)),
			),
		},
	}

	app := &cli.Command{
		Name:    "onair",
		Version: buildVersion,
		Usage:   "monitors your audio/video devices and controls a light based on their status",
		Flags:   globalFlags,
		Commands: []*cli.Command{
			{
				Name:  "hue",
				Usage: "commands for configuring your Hue system",
				Flags: hueFlags,
				Commands: []*cli.Command{
					{
						Name:  "init",
						Usage: "Set up onair for the first time",
						Action: func(ctx context.Context, c *cli.Command) error {
							initHue(ctx, c)
							return nil
						},
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "timeout",
								Usage: "how long to wait for you to push the button on your Hue bridge",
								Value: 30,
							},
						},
					},
					{
						Name:  "lights",
						Usage: "display information about available lights and their color settings",
						Action: func(ctx context.Context, c *cli.Command) error {
							getHueLights(ctx, c)
							return nil
						},
					},
				},
			},
			{
				Name:  "run",
				Usage: "run the log watcher/set your lights",
				Action: func(ctx context.Context, c *cli.Command) error {
					run(ctx, c)
					return nil
				},
				Flags: runFlags,
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("oh no: %s\n", err)
	}
}
