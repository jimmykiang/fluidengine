package main

// Sphere2 is a 2-D sphere geometry.
// Represents 2-D sphere geometry which extends Surface2 by
// overriding surface-related queries.
type Sphere2 struct {
	// Center of the sphere.
	center *Vector3D

	//Radius of the sphere.
	radius float64

	// Local-to-world transform.
	transform *Transform2

	// Flips normal when calling Surface3::closestNormal(...).
	isNormalFlipped bool
}

func NewSphere2(center *Vector3D, radius float64) *Sphere2 {
	return &Sphere2{
		center:          center,
		radius:          radius,
		transform:       NewTransform2(),
		isNormalFlipped: false,
	}
}

// isBounded returns true if bounding box can be defined.
func (s *Sphere2) isBounded() bool {

	return true
}
