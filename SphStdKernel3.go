package main

import "jimmykiang/fluidengine/constants"

// SphStdKernel3 is a standard 3-D SPH kernel function object.
type SphStdKernel3 struct {

	// Kernel radius.
	h float64
	// Square of the kernel radius.
	h2 float64
	// Cubic of the kernel radius.
	h3 float64
	// Fifth-power of the kernel radius.
	h5 float64
}

func NewSphStdKernel3(kernelRadius float64) *SphStdKernel3 {
	h := kernelRadius
	h2 := h * h
	h3 := h2 * h
	h5 := h2 * h3

	return &SphStdKernel3{
		h:  h,
		h2: h2,
		h3: h3,
		h5: h5,
	}
}

// Returns kernel function value at given distance.
func (s *SphStdKernel3) operatorKernel(distance float64) float64 {
	distanceSquared := distance * distance

	if distanceSquared >= s.h2 {
		return 0.0
	} else {
		x := 1 - distanceSquared/s.h2
		return 315.0 / (64 * constants.KPiD * s.h3) * x * x * x
	}
}
