package main

import (
	"flag"
	"hardware-led-cube/frames"
	"hardware-led-cube/ledcube"
	"log"
	"os"
	"time"
)

const FRAMERATE = 60

type Cube interface {
	Render() error
	SetLeds(f frames.Frame) error
	Fini()
}

func main() {
	flag.Parse()
	var cube Cube
	var fs frames.FrameSource

	cube = ledcube.InitCube()
	defer cube.Fini()

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
			if ledcube.Recoverable(err) {
				log.Println(err)
				continue
			}
			break
		}
	}

}
