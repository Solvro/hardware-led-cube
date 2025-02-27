package ledcube

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"hardware-led-cube/frames"
	"log"
)

const (
	WIDTH           = 8
	HEIGHT          = 8
	DEPTH           = 8
	LED_COUNT_HALF  = WIDTH * HEIGHT * DEPTH / 2
	BOTTOM          = 0
	TOP             = 1
	GPIO_PIN_BOTTOM = 18
	GPIO_PIN_TOP    = 19
	BRIGHTNESS      = 64
	FREQ            = 800_000
)

type LedCube ws2811.WS2811

func (c *LedCube) Render() error {
	return (*ws2811.WS2811)(c).Render()
}

func (c *LedCube) SetLeds(f frames.Frame) {
	bottom, top := formatFrame(f.ToXYZ())
	go func() {
		if err := (*ws2811.WS2811)(c).SetLedsSync(BOTTOM, bottom); err != nil {
			log.Panic(err)
		}
	}()
	if err := (*ws2811.WS2811)(c).SetLedsSync(TOP, top); err != nil {
		log.Panic(err)
	}
}

func (c *LedCube) Fini() {
	(*ws2811.WS2811)(c).Fini()
}

func InitCube() *LedCube {
	opt := ws2811.DefaultOptions
	opt.Channels[BOTTOM].GpioPin = GPIO_PIN_BOTTOM
	opt.Channels[TOP].GpioPin = GPIO_PIN_TOP
	opt.Channels[BOTTOM].LedCount = LED_COUNT_HALF
	opt.Channels[TOP].LedCount = LED_COUNT_HALF
	opt.Channels[BOTTOM].Brightness = BRIGHTNESS
	opt.Channels[TOP].Brightness = BRIGHTNESS
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

// Recoverable determines if an error returned by one of the cube's methods is Recoverable
func Recoverable(err error) bool {
	return true
}
