package main

// ImplicitSurfaceSet3 represents 3-D implicit surface set.
type ImplicitSurfaceSet3 struct {
	ImplicitSurface3  *ImplicitSurface3
	bvhInvalidated    bool
	surfaces          []ImplicitSurface3
	unboundedSurfaces []ImplicitSurface3
	bvh               *Bvh3
}

func NewImplicitSurfaceSet3() *ImplicitSurfaceSet3 {
	return &ImplicitSurfaceSet3{
		bvhInvalidated:    true,
		surfaces:          make([]ImplicitSurface3, 0),
		unboundedSurfaces: make([]ImplicitSurface3, 0),
		bvh:               NewBvh3(),
	}
}

// Adds an explicit surface instance.
func (s *ImplicitSurfaceSet3) addExplicitSurface(surface ImplicitSurface3) {

	s.surfaces = append(s.surfaces, surface)
	if !surface.isBounded() {

		s.unboundedSurfaces = append(s.unboundedSurfaces, surface)
	}
	s.invalidateBvh()
}

func (s *ImplicitSurfaceSet3) invalidateBvh() {

	s.bvhInvalidated = true
}
