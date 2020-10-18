package main

import "math"

// Animation represents the base interface for the animation logic in its base level.
type Animation interface {

	// onUpdate should be overriden by downstream structs and implement its logic for updating the animation state.
	onUpdate(*Frame)
}

// SineAnimation contains the evaluated value for a typical sinusoid.
type SineAnimation struct {
	value float64
}

// NewSineAnimation creates and returns a new SineAnimation reference.
func NewSineAnimation() *SineAnimation {
	sineAnimation := &SineAnimation{
		value: 0,
	}

	return sineAnimation
}

// onUpdate for a standard sinusoidal function.
func (sineAnimation *SineAnimation) onUpdate(frame *Frame) {

	sineAnimation.value = math.Sin(10.0 * frame.timeInSeconds())
}
