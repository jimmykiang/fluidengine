package main

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
		return 60 / (kPiD * s.h4) * x
	}
}
