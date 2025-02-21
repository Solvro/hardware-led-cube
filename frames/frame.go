package frames

type FrameSource interface {
	NextFrame() Frame
}

type Frame interface {
	// ToXYZ returns the frame as a 3D array of uint32 colors [x][y][z]uint32 where the uint32 represents an rgb color created like so: 0xRRGGBB
	ToXYZ() [][][]uint32
}
