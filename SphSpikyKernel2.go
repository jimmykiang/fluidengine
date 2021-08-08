package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
)

// Spiky 2-D SPH kernel function object.
// Müller, Matthias, David Charypar, and Markus Gross.
//     "Particle-based fluid simulation for interactive applications."
//     Proceedings of the 2003 ACM SIGGRAPH/Eurographics symposium on Computer
//     animation. Eurographics Association, 2003.
type SphSpikyKernel2 struct {

	// Kernel radius.
	h float64
	// Square of the kernel radius.
	h2 float64
	// Cubic of the kernel radius.
	h3 float64
	// Fourth-power of the kernel radius.
	h4 float64
	// Fifth-power of the kernel radius.
	h5 float64
}

func NewSphSpikyKernel2(h_ float64) *SphSpikyKernel2 {
	h := h_
	h2 := h * h
	h3 := h2 * h
	h4 := h2 * h2
	h5 := h3 * h2

	return &SphSpikyKernel2{
		h:  h,
		h2: h2,
		h3: h3,
		h4: h4,
		h5: h5,
	}
}

func (s *SphSpikyKernel2) secondDerivative(distance float64) float64 {
	if distance >= s.h {
		return 0
	} else {
		x := 1 - distance/s.h
		return 60 / (constants.KPiD * s.h4) * x
	}
}

func (s *SphSpikyKernel2) gradient(
	distance float64,
	directionToCenter *Vector3D.Vector3D,
) *Vector3D.Vector3D {

	a := -s.firstDerivative(distance)
	return directionToCenter.Multiply(a)
}

func (s *SphSpikyKernel2) firstDerivative(distance float64) float64 {

	if distance >= s.h {
		return 0
	} else {
		x := 1.0 - distance/s.h
		return -30.0 / (constants.KPiD * s.h3) * x * x
	}
}

func (s *SphSpikyKernel2) operatorKernel(distance float64) float64 {

	if distance >= s.h {
		return 0.0
	} else {
		x := 1 - distance/s.h
		return 10.0 / (constants.KPiD * s.h2) * x * x * x
	}
}
