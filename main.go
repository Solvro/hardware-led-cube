package main

import (
	"errors"
	"flag"
	"hardware-led-cube/frames"
	"hardware-led-cube/mock"
	"log"
	"os"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	WIDTH           = 4
	HEIGHT          = 4
	DEPTH           = 4
	LED_COUNT_HALF  = WIDTH * HEIGHT * DEPTH / 2
	BOTTOM          = 0
	TOP             = 1
	GPIO_PIN_BOTTOM = 18
	GPIO_PIN_TOP    = 19
	BRIGHTNESS      = 64
	FRAMERATE       = 4
	FREQ            = 800_000
)

var (
	// Command line flags
	is_mock = flag.Bool("mock", true, "Render animation in a window")
)

type Cube interface {
	Render() error
	SetLeds(f frames.Frame)
	Fini()
}

type LedCube ws2811.WS2811

func (c *LedCube) Render() error {
	return (*ws2811.WS2811)(c).Render()
}

func (c *LedCube) SetLeds(f frames.Frame) {
	bottom, top := formatFrame(f.ToXYZ())
	for i := range LED_COUNT_HALF {
		(*ws2811.WS2811)(c).Leds(BOTTOM)[i] = bottom[i]
		(*ws2811.WS2811)(c).Leds(TOP)[i] = top[i]
	}
}

func (c *LedCube) Fini() {
	(*ws2811.WS2811)(c).Fini()
}

func InitLedCube() *LedCube {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = GPIO_PIN_BOTTOM
	opt.Channels[1].GpioPin = GPIO_PIN_TOP
	opt.Channels[0].LedCount = LED_COUNT_HALF
	opt.Channels[1].LedCount = LED_COUNT_HALF
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
	return (*LedCube)(cube)
}

func formatFrame(frame [][][]uint32) (bottom, top [LED_COUNT_HALF]uint32) {
	// Format the frame to match the LED layout
	// TODO: actual implementation
	return bottom, top
}

func main() {
	flag.Parse()
	var cube Cube
	var fs frames.FrameSource

	if !*is_mock {
		cube = InitLedCube()
		defer cube.Fini()
	} else {
		cube = mock.InitCube(WIDTH, HEIGHT, DEPTH)
		defer cube.Fini()
	}

	file, err := os.Open("animation.json")
	if err != nil {
		panic(err)
	}

	fs, ec := frames.NewJSONFileAnimation(file)
	checkParsingError := frames.ErrChanChecker(fs, ec)
	checkParsingError()

	tick := time.Tick(time.Second / FRAMERATE)
	for {
		f := fs.NextFrame()
		checkParsingError()
		cube.SetLeds(f)

		<-tick
		if err := cube.Render(); err != nil {
			if !errors.Is(err, mock.ErrWindowShouldClose) {
				log.Println(err)
				continue
			}
			break
		}
	}

}
