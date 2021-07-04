package main

//type ImplicitSurface2 Plane2D + Sphere2D

type ImplicitSurface2 interface {
	isBounded() bool
	boundingBox() *BoundingBox2D
	signedDistance(otherPoint *Vector3D) float64
}

// Plane2D defines a simple Plane2D struct data representing a 3-D plane geometry.
type Plane2D struct {
	// Base struct for 2-D surface.
	surface2 *Surface2

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
		surface2:        NewSurface2(),
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

// boundingBox returns the bounding box of this surface object.
func (s *Plane2D) boundingBox() *BoundingBox2D {

	return NewBoundingBox2D(NewVector(0, 0, 0), NewVector(0, 0, 0))
}

// Returns the closest point from the given point otherPoint to the surface.
func (p *Plane2D) closestPoint(otherPoint *Vector3D) *Vector3D {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.

	otherPointLocal := p.transform.toLocal(otherPoint)
	d := p.closestPointLocal(otherPointLocal)
	return p.transform.toWorldPointInLocal(d)
}

func (p *Plane2D) closestPointLocal(otherPoint *Vector3D) *Vector3D {

	r := otherPoint.Substract(p.point)
	a := p.normal.DotProduct(r)
	b := p.normal.Multiply(a)
	c := r.Substract(b)
	return c.Add(p.point)

}

func (p *Plane2D) signedDistance(otherPoint *Vector3D) float64 {

	x := p.closestPoint(p.transform.toLocal(otherPoint))
	inside := p.isInside(otherPoint)

	sd := 0.0
	if inside {
		sd = -x.distanceTo(otherPoint)
	} else {
		sd = x.distanceTo(otherPoint)
	}

	if p.isNormalFlipped {
		sd = -sd
	}

	return sd
}

// Returns true if otherPoint is inside the volume defined by the surface.
func (p *Plane2D) isInside(otherPoint *Vector3D) bool {

	return p.isNormalFlipped == !p.isInsideLocal(p.transform.toLocal(otherPoint))
}

// Returns true if otherPoint is inside by given depth the volume
// defined by the surface in local frame.
func (p *Plane2D) isInsideLocal(otherPointLocal *Vector3D) bool {

	cpLocal := p.closestPointLocal(otherPointLocal)
	normalLocal := p.closestNormalLocal(otherPointLocal)
	r := otherPointLocal.Substract(cpLocal)
	return r.DotProduct(normalLocal) < 0.0
}

func (p *Plane2D) closestNormalLocal(otherPoint *Vector3D) *Vector3D {

	return p.normal
}
