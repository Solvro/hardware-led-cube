package ledcube

import (
	"hardware-led-cube/frames"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	WIDTH      = 4
	HEIGHT     = 4
	DEPTH      = 4
	LED_COUNT  = WIDTH * HEIGHT * DEPTH
	CHANNEL    = 0
	GPIO_PIN   = 18
	BRIGHTNESS = 64
	FREQ       = 800_000
)

type LedCube ws2811.WS2811

func (c *LedCube) Render() error {
	return (*ws2811.WS2811)(c).Render()
}

func (c *LedCube) SetLeds(f frames.Frame) error {
	leds_fixed := formatFrame(f.ToXYZ())
	leds := make([]uint32, len(leds_fixed))
	copy(leds, leds_fixed[:])
	return (*ws2811.WS2811)(c).SetLedsSync(CHANNEL, leds)
}

func (c *LedCube) Fini() {
	(*ws2811.WS2811)(c).Fini()
}

func InitCube() *LedCube {
	opt := ws2811.DefaultOptions
	opt.Channels[CHANNEL].GpioPin = GPIO_PIN
	opt.Channels[CHANNEL].LedCount = LED_COUNT
	opt.Channels[CHANNEL].Brightness = BRIGHTNESS
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

func formatFrame(frame [][][]uint32) [LED_COUNT]uint32 {
	// TODO: implement this
	var leds [LED_COUNT]uint32
	return leds
}

// Recoverable determines if an error returned by one of the cube's methods is Recoverable
func Recoverable(err error) bool {
	// TODO: implement this
	return false
}
