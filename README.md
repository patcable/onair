# OnAir

Are you working from home more often in these unprecdented times? Maybe it would be nice to give your housemates an indication that you're on a video or audio call. OnAir will monitor the state of your audio/video input devices, then change the color of a Philips Hue light! It's pretty neat.

To get started you'll need:
* A computer running macOS with [Micro Snitch](https://obdev.at/products/microsnitch/index.html) installed.
* Familiarity with the macOS command line interface.
* A [Philips Hue Bridge](https://www2.meethue.com/en-us/p/hue-bridge/046677458478).
* Some form of light to use with it - perhaps a [Philips Hue Go portable light](https://www2.meethue.com/en-us/p/hue-white-and-color-ambiance-go-portable-light/714606048).

## Building OnAir
Never used Go before? You can use [homebrew](https://brew.sh/) to install golang and get _going_ pretty quickly:

```
brew install golang
mkdir ~/gostuff
export GOPATH=$HOME/gostuff
export PATH=$GOPATH/bin:$PATH
go install github.com/patcable/onair
```

You can put your $GOPATH wherever you want. You'll likely want to toss those `export` lines in your shell profile. If you want to modify any of the code, you can find that in `$GOPATH/src/github.com/patcable/onair`.

## Configuring OnAir
### Creating a Username on your Hue bridge
Once you have the `onair` binary, you'll need to set OnAir up so that it can talk to the Hue bridge. To do this, 
run `onair init`, walk over to the Hue bridge and press the button. This authenticates OnAir to the Hue bridge.

If successful, you should see this output:

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
Now you can run `onair lights` which will display information about the lights connected to the bridge. The output should look like:
```
Your Hue bridge knows about these lights:
 -> hue boi:
      ID:         5
      Color Mode: xy
      Brightness: 70
      Hue:        34758
      Saturation: 212
      Color Temp: 153
      XY:         [0.1988 0.4146]
```

The name of the light is to the right of the arrow. Take note of the ID number. Figure out what color you want to use to indicate that your camera or microphone is active. I use an orange, something like `0.5153,0.4146`. Then figure out what color you want to use to indicate inactivity. I have a mint green - `0.1989,0.4146`. You can use different brightnesses too, from 1 to 254. I use `70`.

Now that you have that information, update your YAML file:
```
---
hueuid: wrheWt05IVRdftSh76Af8fAUcX0o4Olwi-YvNYIO
light: 5
brightness: 70
active: 0.5153,0.4146
inactive: 0.1989,0.4146
```
### Run the watcher!
At this point, you're ready to give it a go. Make sure that Micro Snitch is running and type `onair run` then hit enter. You'll see a log line indicating that OnAir is watching the Micro Snitch log, but other than that it will be quiet:
```
2020/04/18 20:16:51 Seeked /Users/cable/Library/Logs/Micro Snitch.log - &{Offset:0 Whence:2}
```

The lack of output is normal. Find an application that enables the microphone or camera, and watch your Hue change color! It should change color back to the inactive color once you close the application. This works across any application on macOS!

## Contributing, Questions, etc.
PRs are welcome! I'm happy to look at those as time allows. If you run into issues, opening an issue is the best way to get my attention, but it might be a while.

If you're interested in contributing, the application is pretty small. 
* `onair.go` has all the configurable variables that are used throughout the application
* `watcher.go` contains the function for 
* `hue.go` contains functions for communicating with the Hue bridge.

## To Do
* Maybe a debug mode would be nice
* I built this for Micro Snitch because that's what I have, but extending it for [Oversight](https://objective-see.com/products/oversight.html)
  wouldn't be too bad (likely an update to the grok pattern in `watcher.go` and some variable name changes)
* Write a launch agent definition for this so that it can load on user login
* Suggestions Welcome!
