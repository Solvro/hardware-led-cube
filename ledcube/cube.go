//go:build !arm64 || !linux

package ledcube

import (
	"errors"
	"hardware-led-cube/mock"
)

const (
	WIDTH  = 4
	HEIGHT = 4
	DEPTH  = 4
)

func InitCube() *mock.Cube {
	return mock.InitCube(WIDTH, HEIGHT, DEPTH)
}

// Recoverable determines if an error returned by one of the cube's methods is Recoverable
func Recoverable(err error) bool {
	return !errors.Is(err, mock.ErrWindowShouldClose)
}
