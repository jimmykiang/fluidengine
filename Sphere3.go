package main

import "jimmykiang/fluidengine/Vector3D"

// Sphere3 is a 3-D sphere geometry.
// Represents 3-D sphere geometry which extends Surface3 by
// overriding surface-related queries.
type Sphere3 struct {

	// Base struct for 3-D surface.
	surface3 *Surface3

	// Center of the sphere.
	center *Vector3D.Vector3D

	//Radius of the sphere.
	radius float64

	// Local-to-world transform.
	transform *Transform3

	// Flips normal when calling Surface3::closestNormal(...).
	isNormalFlipped bool
}

func NewSphere3(center *Vector3D.Vector3D, radius float64) *Sphere3 {
	return &Sphere3{
		surface3:        NewSurface3(),
		center:          center,
		radius:          radius,
		transform:       NewTransform3(),
		isNormalFlipped: false,
	}
}

// isBounded returns true if bounding box can be defined.
func (s *Sphere3) isBounded() bool {

	return true
}

// boundingBox returns the bounding box of this surface object.
func (s *Sphere3) boundingBox() *BoundingBox3D {

	r := Vector3D.NewVector(s.radius, s.radius, s.radius)
	return s.transform.toWorldBoundingBox(NewBoundingBox3D(s.center.Substract(r), s.center.Add(r)))
}

func (s *Sphere3) signedDistance(otherPoint *Vector3D.Vector3D) float64 {

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

// Returns the closest point from the given point otherPoint to the surface.
func (s *Sphere3) closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.

	otherPointLocal := s.transform.toLocal(otherPoint)
	d := s.closestPointLocal(otherPointLocal)
	return s.transform.toWorld(d)
}

// Returns true if otherPoint is inside the volume defined by the surface.
func (s *Sphere3) isInside(otherPoint *Vector3D.Vector3D) bool {

	return s.isNormalFlipped == !s.isInsideLocal(s.transform.toLocal(otherPoint))
}

func (s *Sphere3) closestPointLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := s.closestNormalLocal(otherPoint).Multiply(s.radius)

	return r.Add(s.center)

}

func (s *Sphere3) closestNormalLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	if s.center.IsSimilar(otherPoint) {
		return Vector3D.NewVector(1, 0, 0)
	} else {

		return otherPoint.Substract(s.center).Normalize()
	}
}

// Returns true if otherPoint is inside by given depth the volume
// defined by the surface in local frame.
func (s *Sphere3) isInsideLocal(otherPointLocal *Vector3D.Vector3D) bool {

	cpLocal := s.closestPointLocal(otherPointLocal)
	normalLocal := s.closestNormalLocal(otherPointLocal)
	r := otherPointLocal.Substract(cpLocal)
	return r.DotProduct(normalLocal) < 0.0
}

func (s *Sphere3) getTransform() *Transform3 {
	return s.transform
}
