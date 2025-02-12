package main

import (
	"encoding/json"
	"io"
)

type JSONFileAnimation struct {
	frames    []frame
	frameIdx  int
	frameChan chan frame
}

type frame [][][]uint32

func (f frame) Normalize() (bottom [LED_COUNT]uint32, top [LED_COUNT]uint32) {
	// TODO: find out what the actual output format needed is, and reimplement accordingly
	for i := 0; i < LED_COUNT; i++ {
		bottom[i] = f[i%WIDTH][(i/HEIGHT)%(HEIGHT/2)][i/(WIDTH*HEIGHT/2)]
		top[i] = f[i%WIDTH][(i/HEIGHT)%(HEIGHT/2)+HEIGHT/2][i/(WIDTH*HEIGHT/2)]
	}
	return
}

func NewJSONFileAnimation(input io.ReadCloser) (jfa *JSONFileAnimation, ec chan error) {
	ec = make(chan error, 1)
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

	jfa = &JSONFileAnimation{frames: make([]frame, 0), frameChan: make(chan frame, 32)}
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

func (jfa *JSONFileAnimation) NextFrame() Frame {
	frame, ok := <-jfa.frameChan
	if ok {
		return frame
	}

	frame = jfa.frames[jfa.frameIdx]
	jfa.frameIdx = (jfa.frameIdx + 1) % len(jfa.frames)
	return frame
}
