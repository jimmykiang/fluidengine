package main

// SurfaceToImplicit2 is a implicit surface wrapper for generic Surface2 instance.
// This struct represents 2-D implicit surface that converts Surface2 instance
// to an ImplicitSurface2 object. The conversion is made by evaluating closest
// point and normal from a given point for the given (explicit) surface. Thus,
// this conversion won't work for every single surfaces. Use this class only
// for the basic primitives such as Sphere2 or Box2.
type SurfaceToImplicit2 struct {

	// Not Needed?
	surface *Plane2D
}

func NewSurfaceToImplicit2(surface *Plane2D) *SurfaceToImplicit2 {
	return &SurfaceToImplicit2{surface: surface}
}
