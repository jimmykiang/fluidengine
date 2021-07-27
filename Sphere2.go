package main

import "jimmykiang/fluidengine/Vector3D"

// Sphere2 is a 2-D sphere geometry.
// Represents 2-D sphere geometry which extends Surface2 by
// overriding surface-related queries.
type Sphere2 struct {

	// Base struct for 2-D surface.
	surface2 *Surface2

	// Center of the sphere.
	center *Vector3D.Vector3D

	//Radius of the sphere.
	radius float64

	// Local-to-world transform.
	transform *Transform2

	// Flips normal when calling Surface3::closestNormal(...).
	isNormalFlipped bool
}

func NewSphere2(center *Vector3D.Vector3D, radius float64) *Sphere2 {
	return &Sphere2{
		surface2:        NewSurface2(),
		center:          center,
		radius:          radius,
		transform:       NewTransform2(),
		isNormalFlipped: false,
	}
}

func (s *Sphere2) getTransform() *Transform2 {
	return s.transform
}

// isBounded returns true if bounding box can be defined.
func (s *Sphere2) isBounded() bool {

	return true
}

// boundingBox returns the bounding box of this surface object.
func (s *Sphere2) boundingBox() *BoundingBox2D {

	r := Vector3D.NewVector(s.radius, s.radius, 0)
	return s.transform.toWorld(NewBoundingBox2D(s.center.Substract(r), s.center.Add(r)))
}

// Returns true if otherPoint is inside by given depth the volume
// defined by the surface in local frame.
func (s *Sphere2) isInsideLocal(otherPointLocal *Vector3D.Vector3D) bool {

	cpLocal := s.closestPointLocal(otherPointLocal)
	normalLocal := s.closestNormalLocal(otherPointLocal)
	r := otherPointLocal.Substract(cpLocal)
	return r.DotProduct(normalLocal) < 0.0
}

// Returns the closest point from the given point otherPoint to the surface.
func (s *Sphere2) closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.

	otherPointLocal := s.transform.toLocal(otherPoint)
	d := s.closestPointLocal(otherPointLocal)
	return s.transform.toWorldPointInLocal(d)
}

func (s *Sphere2) closestNormalLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	if s.center.IsSimilar(otherPoint) {
		return Vector3D.NewVector(1, 0, 0)
	} else {

		return otherPoint.Substract(s.center).Normalize()
	}

}

func (s *Sphere2) closestPointLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := s.closestNormalLocal(otherPoint).Multiply(s.radius)

	return r.Add(s.center)

}

// Returns true if otherPoint is inside the volume defined by the surface.
func (s *Sphere2) isInside(otherPoint *Vector3D.Vector3D) bool {

	return s.isNormalFlipped == !s.isInsideLocal(s.transform.toLocal(otherPoint))
}

func (s *Sphere2) signedDistance(otherPoint *Vector3D.Vector3D) float64 {

	x := s.closestPoint(s.transform.toLocal(otherPoint))

	inside := s.isInside(otherPoint)

	sd := 0.0
	if inside {
		sd = -x.DistanceTo(otherPoint)
	} else {
		sd = x.DistanceTo(otherPoint)
	}
	if s.isNormalFlipped {
		sd = -sd
	}
	return sd
}
