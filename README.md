# OnAir

Are you working from home more often in these unprecedented times? Maybe it would be nice to give your housemates an indication that you're on a video or audio call. OnAir will monitor the state of your audio/video input devices, then change the color of a light! Here's an animated gif of it doing the thing:

![onair](https://user-images.githubusercontent.com/2254952/79677012-47662c80-81ba-11ea-966b-99fd86452e41.gif)

## Neat I want one!

You'll need the following to get started:
* A computer running macOS with [Micro Snitch](https://obdev.at/products/microsnitch/index.html) installed.
* Familiarity with the macOS command line interface.
* A [Philips Hue Bridge](https://www2.meethue.com/en-us/p/hue-bridge/046677458478).
* Some form of light to use with it - perhaps a [Philips Hue Go portable light](https://www2.meethue.com/en-us/p/hue-white-and-color-ambiance-go-portable-light/714606048).

### Building OnAir
Never used Go before? You can use [homebrew](https://brew.sh/) to install golang and get _going_ pretty quickly:

```
brew install golang
mkdir ~/gostuff
export GOPATH=$HOME/gostuff
export PATH=$GOPATH/bin:$PATH
go get github.com/patcable/onair
go install github.com/patcable/onair
```

You can put your $GOPATH wherever you want. You'll likely want to toss those `export` lines in your shell profile. If you want to modify any of the code, you can find that in `$GOPATH/src/github.com/patcable/onair`.

## Configuring OnAir

This guide focuses on a Philips Hue system; if other lights become supported I'll figure out what to do with this guide.

### Creating a Username on your Hue bridge
Once you have the `onair` binary, you'll need to set OnAir up so that it can talk to the Hue bridge. To do this, run `onair hue init`, walk over to the Hue bridge and press the button. This authenticates OnAir to the Hue bridge. If successful, you should see this output:

```
Head on over to your HUE Bridge and push the button to create a new user.
Waiting 15 seconds - you can set this longer by passing --timeout to 'onair init'.
   Username: wrheWt05IVRdftSh76Af8fAUcX0o4Olwi-YvNYIO
Save this value, you'll need it for the config file.
```

Now, create a file at `$HOME/.onair.yml` with this in it:

```
---
hueuid: wrheWt05IVRdftSh76Af8fAUcX0o4Olwi-YvNYIO
```

### Figuring out your lighting situation
Once the `hueuid` value is in place, you can run `onair hue lights` which will display information about the lights connected to the bridge. The output should look like:

```
Your Hue bridge knows about these lights:
 -> Meeting Orb:
      ID:         5
      Color Mode: xy
      Brightness: 70
      Hue:        34758
      Saturation: 212
      Color Temp: 153
      XY:         [0.1988 0.4146]
```

The name of the light is to the right of the arrow. Take note of the `ID` number. Then, figure out what color you want to use to indicate that your camera or microphone is active. Use the Hue application on your phone to change the color, then run `onair hue lights` to look at the new settings.

I use an orange, something like `0.5153,0.4146` for an active call, and a mint green - `0.1989,0.4146` for an inactive call. You can use different brightnesses too, from 1 to 254. I use `70`.

Now that you've figured out what colors you want to use, update your YAML file:

```
---
hueuid: wrheWt05IVRdftSh76Af8fAUcX0o4Olwi-YvNYIO
huelight: 5
huebrightness: 70
hueactive: 0.5153,0.4146
hueinactive: 0.1989,0.4146
```

### Run the watcher!
At this point, you're ready to give it a go. Make sure that Micro Snitch is running and type `onair run` then hit enter. You'll see a log line indicating that OnAir is watching the Micro Snitch log, but other than that it will be quiet:
```
2020/04/18 20:16:51 Seeked /Users/cable/Library/Logs/Micro Snitch.log - &{Offset:0 Whence:2}
```

The lack of output is normal. Find an application that enables the microphone or camera, and watch your light change color! It should change color back to the inactive color once you close the application. This works across any application on macOS! It's very satisfying.

## Issues
PRs are welcome! I'm happy to look at those as time allows. If you run into issues, opening an issue is the best way to get my attention, but it might be a while. I should probably add in better debug log support to make this a bit easier on myself. Coming soon :)

## Contributing

If you're interested in adding bits to OnAir this section will be helpful.

### App Layout
* `onair.go` has the CLI skeleton, and all the configuration options. If you want to add a config option or change a default that isn't in the config file, this is where you'd do it. Any command that has a `Before` attribute that references the `altsrc` package can be pulled in from the YAML file referenced in the `config` CLI String Flag.
* `watcher.go` has the main run function of the application. This is the part of the app that actually watches the log file you specify, configures the parser for that file, sets up the lighting options based on configured system and type, etc. 
* `hue.go` has the functions specific to the hue lighting setup. Notably, a `hueConfig` struct that gets used in `watcher.go` to actually set stuff for the lights.

### Adding a log parser
If you want to add a parser, update `configureParser` with a new grok pattern (`gr.AddPattern`) and then update the switch statement in `configureParser` so that the run function knows to use that grok filter. You'll need to update the for loop in the run function as well since there are likely new field names that you may need to use.

### Adding a lighting system
If you'd like to add a different lighting system, you should check to see if it has a golang library first. If not, step one may be: write a library for your lighting system. Once that's done, to get that into OnAir, you'll need to:

* Add some commands/subcommands to get information about your lights in `onair.go`. Config values should be named `SYSTEMthing` (ie. `lifxlight` if you were adding in lifx support)
* Write those functions in `SYSTEM.go` (no caps plz) 
* Update `configureLightSystem` and `setLight` functions

### Improvement Ideas
* Probably better error handling or flow. Using logrus or having some sort of debug mode might be nice.
* I built this for Micro Snitch because that's what I have, but extending it for [Oversight](https://objective-see.com/products/oversight.html) should be possible depending on the log type.
* Write a launch agent definition for this so that it can load on user login
