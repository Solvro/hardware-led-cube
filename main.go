package main

import (
	"flag"
	"fmt"
	mock "hardware-led-cube/mock"
	"log"
	"os"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
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
	mock_bool = flag.Bool("mock", false, "Render animation in a window")
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

type mockCube struct {
}

func (c *mockCube) Render() error {
	return nil
}

func (c *mockCube) SetLeds(f Frame) {
}

func (c *mockCube) Fini() {
}

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

func InitMockCube() Cube {
	return &mockCube{}
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
	ledUpdates := make(chan [][][]mock.Led)

	if !*mock_bool {
		cube = InitLedCube()
		defer cube.Fini()
	} else {
		go mock.Main(ledUpdates)
		cube = InitMockCube()
	}

	file, err := os.Open("animation.json")
	if err != nil {
		panic(err)
	}
	fs, ec := NewJSONFileAnimation(file)
	log.Println(file.Name())
	checkParsingError := errChanChecker(fs, ec)
	checkParsingError()
	tick := time.Tick(time.Second / FRAMERATE)

	for {
		f := fs.NextFrame()
		checkParsingError()
		cube.SetLeds(f)

		if *mock_bool {
			bottom, top := f.Normalize()
			leds := make([][][]mock.Led, WIDTH)
			for x := 0; x < WIDTH; x++ {
				leds[x] = make([][]mock.Led, HEIGHT)
				for y := 0; y < HEIGHT; y++ {
					leds[x][y] = make([]mock.Led, DEPTH)
					for z := 0; z < DEPTH; z++ {
						leds[x][y][z] = mock.Led{
							X: float32(x),
							Y: float32(y),
							Z: float32(z),
							R: float32(bottom[x*HEIGHT+y]) / 255.0,
							G: float32(top[x*HEIGHT+y]) / 255.0,
							B: float32(bottom[x*HEIGHT+y]) / 255.0,
						}
					}
				}
			}
			ledUpdates <- leds
		}

		<-tick
		err := cube.Render()
		if err != nil {
			log.Println(err)
		}
	}
}
