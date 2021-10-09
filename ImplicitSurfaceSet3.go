package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"math"
)

// ImplicitSurfaceSet3 represents 3-D implicit surface set.
type ImplicitSurfaceSet3 struct {
	ImplicitSurface3  *ImplicitSurface3
	bvhInvalidated    bool
	surfaces          []ImplicitSurface3
	unboundedSurfaces []ImplicitSurface3
	bvh               *Bvh3
	transform         *Transform3
	isNormalFlipped   bool
}

func NewImplicitSurfaceSet3() *ImplicitSurfaceSet3 {
	return &ImplicitSurfaceSet3{
		bvhInvalidated:    true,
		surfaces:          make([]ImplicitSurface3, 0),
		unboundedSurfaces: make([]ImplicitSurface3, 0),
		bvh:               NewBvh3(),
		transform:         NewTransform3(),
		isNormalFlipped:   false,
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

// updateQueryEngine updates internal spatial query engine.
func (s *ImplicitSurfaceSet3) updateQueryEngine() {
	s.invalidateBvh()
	s.buildBvh()
}
func (s *ImplicitSurfaceSet3) buildBvh() {
	if s.bvhInvalidated {

		surfs := make([]ImplicitSurface3, 0, 0)
		bounds := make([]*BoundingBox3D, 0, 0)

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

// isBounded returns true if bounding box can be defined.
func (s *ImplicitSurfaceSet3) isBounded() bool {
	// All surfaces should be bounded.
	for _, surface := range s.surfaces {
		if !surface.isBounded() {
			return false
		}
	}

	// Empty set is not bounded.
	return len(s.surfaces) != 0
}

func (s *ImplicitSurfaceSet3) signedDistance(otherPoint *Vector3D.Vector3D) float64 {

	t := s.transform.toLocal(otherPoint)
	sd := s.signedDistanceLocal(t)

	if s.isNormalFlipped {
		sd = -sd
	}
	return sd
}

func (s *ImplicitSurfaceSet3) signedDistanceLocal(otherPoint *Vector3D.Vector3D) float64 {
	sdf := math.MaxFloat64
	for _, surface := range s.surfaces {

		sdf = math.Min(sdf, surface.signedDistance(otherPoint))
	}
	return sdf
}
