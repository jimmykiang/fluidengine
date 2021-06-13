package main

// ImplicitSurfaceSet2 represents 2-D implicit surface set.
type ImplicitSurfaceSet2 struct {
	ImplicitSurface2  *ImplicitSurface2
	bvhInvalidated    bool
	surfaces          []ImplicitSurface2
	unboundedSurfaces []ImplicitSurface2
	bvh               *Bvh2
}

func NewImplicitSurfaceSet2() *ImplicitSurfaceSet2 {
	return &ImplicitSurfaceSet2{
		bvhInvalidated:    true,
		surfaces:          make([]ImplicitSurface2, 0),
		unboundedSurfaces: make([]ImplicitSurface2, 0),
		//bvh,
	}
}

// Adds an explicit surface instance.
func (s *ImplicitSurfaceSet2) addExplicitSurface(surface ImplicitSurface2) {

	s.surfaces = append(s.surfaces, surface)
	if !surface.isBounded() {

		s.unboundedSurfaces = append(s.unboundedSurfaces, surface)
	}
	s.invalidateBvh()
}

func (s *ImplicitSurfaceSet2) invalidateBvh() {

	s.bvhInvalidated = true
}
