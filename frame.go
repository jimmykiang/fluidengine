package main

// Frame contains the representation of an animation frame.
// This struct holds current animation frame index and frame interval in
// seconds.
type Frame struct {
	index                 int
	timeIntervalInSeconds float64
}

// NewFrame creates and returns a new Frame reference.
func NewFrame() *Frame {
	frame := &Frame{
		timeIntervalInSeconds: 1.0 / 60.0,
		index:                 0,
	}

	return frame
}

func (frame *Frame) timeInSeconds() float64 {

	return float64(frame.index) * frame.timeIntervalInSeconds
}

func (frame *Frame) advance() {

	frame.index++
}
