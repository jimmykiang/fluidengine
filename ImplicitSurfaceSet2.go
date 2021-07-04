package main

import "math"

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
		bvh:               NewBvh2(),
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

// updateQueryEngine updates internal spatial query engine.
func (s *ImplicitSurfaceSet2) updateQueryEngine() {
	s.invalidateBvh()
	s.buildBvh()
}

func (s *ImplicitSurfaceSet2) buildBvh() {
	if s.bvhInvalidated {

		surfs := make([]ImplicitSurface2, 0, 0)
		bounds := make([]*BoundingBox2D, 0, 0)

		for i := 0; i < len(s.surfaces); i++ {
			if s.surfaces[i].isBounded() {
				surfs = append(surfs, s.surfaces[i])
				bounds = append(bounds, s.surfaces[i].boundingBox())
			}
		}
		s.bvh.build(surfs, bounds)
		s.bvhInvalidated = false
	}
}

func (s *ImplicitSurfaceSet2) isBounded() bool {
	// All surfaces should be bounded.
	for _, surface := range s.surfaces {
		if !surface.isBounded() {
			return false
		}
	}

	// Empty set is not bounded.
	return len(s.surfaces) != 0
}

func (s *ImplicitSurfaceSet2) signedDistance(candidate *Vector3D) float64 {
	sdf := math.MaxFloat64
	for _, surface := range s.surfaces {

		sdf = math.Min(sdf, surface.signedDistance(candidate))
	}
	return sdf
}
