package frames

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonFileAnimation struct {
	frames    []frame
	frameIdx  int
	frameChan chan frame
}

type frame [][][]uint32

func (f frame) ToXYZ() [][][]uint32 {
	return f
}

// NewJSONFileAnimation creates a new jsonFileAnimation from the given input
// which needs to be a json in the format [frame][x][y][z]uint32 and uint32 represents an rgb color created like so: 0xRRGGBB
func NewJSONFileAnimation(input io.ReadCloser) (*jsonFileAnimation, <-chan error) {
	ec := make(chan error, 1)
	defer input.Close()
	dec := json.NewDecoder(input)

	t, err := dec.Token()
	if err != nil {
		ec <- err
		return nil, ec
	}

	if t != json.Delim('[') {
		ec <- io.ErrUnexpectedEOF
		return nil, ec
	}

	jfa := &jsonFileAnimation{frames: make([]frame, 0), frameChan: make(chan frame, 32)}
	go func(fc chan<- frame, frames *[]frame) {
		defer close(fc)
		for dec.More() {
			var f frame
			if err := dec.Decode(&f); err != nil {
				if err == io.EOF {
					break
				}
				ec <- err
				return
			}
			*frames = append(*frames, f)
			fc <- f
		}
	}(jfa.frameChan, &jfa.frames)
	return jfa, ec
}

func (jfa *jsonFileAnimation) NextFrame() Frame {
	frame, ok := <-jfa.frameChan
	if ok {
		return frame
	}

	frame = jfa.frames[jfa.frameIdx]
	jfa.frameIdx = (jfa.frameIdx + 1) % len(jfa.frames)
	return frame
}

// ErrChanChecker returns a function that checks if the channel containing FrameSource errors has anything in it, and closes the program if it does after printing error info
func ErrChanChecker(fs FrameSource, ec <-chan error) func() {
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
