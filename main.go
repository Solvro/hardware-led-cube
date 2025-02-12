package main

import (
	"flag"
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"log"
	"os"
	"time"
)

const (
	WIDTH           = 8
	HEIGHT          = 8
	DEPTH           = 8
	LED_COUNT       = WIDTH * HEIGHT * DEPTH / 2
	BOTTOM          = 0
	TOP             = 1
	GPIO_PIN_BOTTOM = 18
	GPIO_PIN_TOP    = 19
	BRIGHTNESS      = 64
	FRAMERATE       = 60
	FREQ            = 800_000
)

var (
	// Command line flags
	mock = flag.Bool("mock", false, "Render animation in a window")
)

type FrameSource interface {
	NextFrame() Frame
}

type Frame interface {
	// Normalize the frame to the format the cube expects
	Normalize() (bottom [LED_COUNT]uint32, top [LED_COUNT]uint32)
}

type Cube interface {
	Render() error
	SetLeds(f Frame)
	Fini()
}

type ledCube ws2811.WS2811

func (c *ledCube) Render() error {
	return (*ws2811.WS2811)(c).Render()
}

func (c *ledCube) SetLeds(f Frame) {
	bottom, top := f.Normalize()
	for i := range LED_COUNT {
		(*ws2811.WS2811)(c).Leds(BOTTOM)[i] = bottom[i]
		(*ws2811.WS2811)(c).Leds(TOP)[i] = top[i]
	}
}

func (c *ledCube) Fini() {
	(*ws2811.WS2811)(c).Fini()
}

func InitLedCube() *ledCube {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = GPIO_PIN_BOTTOM
	opt.Channels[1].GpioPin = GPIO_PIN_TOP
	opt.Channels[0].LedCount = LED_COUNT
	opt.Channels[1].LedCount = LED_COUNT
	opt.Channels[0].Brightness = BRIGHTNESS
	opt.Channels[1].Brightness = BRIGHTNESS
	opt.Frequency = FREQ

	cube, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}

	err = cube.Init()
	if err != nil {
		panic(err)
	}
	return (*ledCube)(cube)
}

// returns a function that checks if the channel containing FrameSource errors has anything in it, and closes the program if it does after printing error info
func errChanChecker(fs FrameSource, ec <-chan error) func() {
	return func() {
		select {
		case err := <-ec:
			if err != nil {
				panic(fmt.Sprintf("An error occurred while parsing frames from FrameSource:\n"+
					"%v\nError:\n%v", fs, err))
			}
		default:
		}
	}
}

func main() {
	flag.Parse()
	var cube Cube
	var fs FrameSource
	if !*mock {
		cube = InitLedCube()
		defer cube.Fini()
	} else {
		cube = InitMockCube()
	}

	file, err := os.Open("animation.json")
	if err != nil {
		panic(err)
	}
	fs, ec := NewJSONFileAnimation(file)
	checkParsingError := errChanChecker(fs, ec)
	checkParsingError()
	tick := time.Tick(time.Second / FRAMERATE)

	for {
		f := fs.NextFrame()
		checkParsingError()
		cube.SetLeds(f)

		<-tick
		err := cube.Render()
		if err != nil {
			log.Println(err)
		}
	}
}
