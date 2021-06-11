package main

// SphStdKernel2 is a standard 2-D SPH kernel function object.
type SphStdKernel2 struct {

	// Kernel radius.
	h float64
	// Square of the kernel radius.
	h2 float64
	// Cubic of the kernel radius.
	h3 float64
	// Fourth-power of the kernel radius.
	h4 float64
}

func NewSphStdKernel2(h_ float64) *SphStdKernel2 {
	h := h_
	h2 := h * h
	h3 := h2 * h
	h4 := h2 * h2

	return &SphStdKernel2{
		h:  h,
		h2: h2,
		h3: h3,
		h4: h4,
	}
}

func (s *SphStdKernel2) operatorKernel(distance float64) float64 {
	distanceSquared := distance * distance

	if distanceSquared >= s.h2 {
		return 0.0
	} else {
		x := 1 - distanceSquared/s.h2
		return 4.0 / (kPiD * s.h2) * x * x * x
	}
}
