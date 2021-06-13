package main

//type ImplicitSurface2 Plane2D + Sphere2D

type ImplicitSurface2 interface {
	isBounded() bool
}

// Plane2D defines a simple Plane2D struct data representing a 3-D plane geometry.
type Plane2D struct {
	// Plane normal.
	normal *Vector3D

	// Point that lies on the plane.
	point *Vector3D

	// Local-to-world transform.
	transform *Transform2

	// Flips normal when calling Surface3::closestNormal(...).
	isNormalFlipped bool
}

// NewPlane2D constructs a plane that cross \p point with surface normal \p normal.
func NewPlane2D(normal, point *Vector3D) *Plane2D {
	return &Plane2D{
		normal:          normal,
		point:           point,
		transform:       NewTransform2(),
		isNormalFlipped: false,
	}
}

// isBounded returns true if bounding box can be defined.
func (p *Plane2D) isBounded() bool {

	return false
}
