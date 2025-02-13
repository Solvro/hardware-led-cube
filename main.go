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
	WIDTH           = 4
	HEIGHT          = 4
	DEPTH           = 4
	LED_COUNT_HALF  = WIDTH * HEIGHT * DEPTH / 2
	BOTTOM          = 0
	TOP             = 1
	GPIO_PIN_BOTTOM = 18
	GPIO_PIN_TOP    = 19
	BRIGHTNESS      = 64
	FRAMERATE       = 1
	FREQ            = 800_000
)

var (
	// Command line flags
	mock_bool = flag.Bool("mock", true, "Render animation in a window")
)

type FrameSource interface {
	NextFrame() Frame
}

type Frame interface {
	// Normalize the frame to the format the cube expects
	Normalize() (bottom [LED_COUNT_HALF]uint32, top [LED_COUNT_HALF]uint32)
}

type Cube interface {
	Render() error
	SetLeds(f Frame)
	Fini()
	Run()
}

type ledCube ws2811.WS2811

type mockCube struct {
	leds       [][][]mock.Led
	ledUpdates chan [][][]mock.Led
}

func (c *mockCube) Render() error {
	c.ledUpdates <- c.leds

	return nil
}

func (c *mockCube) SetLeds(f Frame) {
	bottom, top := f.Normalize()
	for x := 0; x < DEPTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			for z := 0; z < WIDTH; z++ {
				var color uint32
				if x < DEPTH/2 {
					color = bottom[x*HEIGHT*WIDTH+y*WIDTH+z]
				} else {
					color = top[(x-DEPTH/2)*HEIGHT*WIDTH+y*WIDTH+z]
				}
				led := c.leds[x][y][z]
				led.R = float32((color >> 16) & 0xFF)
				led.G = float32((color >> 8) & 0xFF)
				led.B = float32(color & 0xFF)
				c.leds[x][y][z] = led
			}
		}
	}
}

func (c *mockCube) Fini() {
}

func (c *mockCube) Run() {
	mock.Setup(c.ledUpdates)
}

func (c *ledCube) Render() error {
	return (*ws2811.WS2811)(c).Render()
}

func (c *ledCube) SetLeds(f Frame) {
	bottom, top := f.Normalize()
	for i := range LED_COUNT_HALF {
		(*ws2811.WS2811)(c).Leds(BOTTOM)[i] = bottom[i]
		(*ws2811.WS2811)(c).Leds(TOP)[i] = top[i]
	}
}

func (c *ledCube) Fini() {
	(*ws2811.WS2811)(c).Fini()
}

func (c *ledCube) Run() {

}

func InitLedCube() *ledCube {
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
	return (*ledCube)(cube)
}

func InitMockCube() Cube {
	leds := make([][][]mock.Led, DEPTH)
	for z := 0; z < DEPTH; z++ {
		leds[z] = make([][]mock.Led, HEIGHT)
		for y := 0; y < HEIGHT; y++ {
			leds[z][y] = make([]mock.Led, WIDTH)
			for x := 0; x < WIDTH; x++ {
				leds[z][y][x] = mock.Led{
					X: float32(x),
					Y: float32(y),
					Z: float32(z),
					R: 1.0,
					G: 1.0,
					B: 1.0}
			}
		}
	}

	ledUpdates := make(chan [][][]mock.Led)

	return &mockCube{
		leds:       leds,
		ledUpdates: ledUpdates,
	}
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

	if !*mock_bool {
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

	go func() {
		for {
			f := fs.NextFrame()
			checkParsingError()
			cube.SetLeds(f)

			<-tick
			log.Println(time.Now())
			if err := cube.Render(); err != nil {
				log.Println(err)
			}
		}
	}()

	if *mock_bool {
		cube.Run()
	} else {
		select {}
	}
}
