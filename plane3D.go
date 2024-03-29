package main

import "jimmykiang/fluidengine/Vector3D"

// Plane3D defines a simple Plane3D struct data representing a 3-D plane geometry.
type Plane3D struct {
	// Plane normal.
	normal *Vector3D.Vector3D

	// Point that lies on the plane.
	point *Vector3D.Vector3D

	// Local-to-world transform.
	transform *Transform3

	// Flips normal when calling Surface3::closestNormal(...).
	isNormalFlipped bool
}

func (p *Plane3D) closestPointLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := otherPoint.Substract(p.point)
	a := p.normal.DotProduct(r)
	b := p.normal.Multiply(a)
	c := r.Substract(b)
	d := c.Add(p.point)
	return d
}

// Returns the normal to the closest point on the surface from the given otherPoint.
func (p *Plane3D) closestDistance(otherPoint *Vector3D.Vector3D) float64 {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.
	otherPointLocal := p.transform.toLocal(otherPoint)
	d := p.closestPointLocal(otherPointLocal)

	//Returns the distance to the other vector.1

	return otherPointLocal.Substract(d).Length()
}

// Returns the closest point from the given point otherPoint to the surface.
func (p *Plane3D) closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.

	otherPointLocal := p.transform.toLocal(otherPoint)
	d := p.closestPointLocal(otherPointLocal)
	return p.transform.toWorld(d)
}

func (p *Plane3D) closestNormal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	result := p.transform.toWorldDirection(p.closestNormalLocal(otherPoint))
	if p.isNormalFlipped {

		result.Multiply(-1)
	}

	return result
}

func (p *Plane3D) closestNormalLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	return p.normal
}

// Returns true if otherPoint is inside the volume defined by the surface.
func (p *Plane3D) isInside(otherPoint *Vector3D.Vector3D) bool {

	return p.isNormalFlipped == !p.isInsideLocal(p.transform.toLocal(otherPoint))
}

// Returns true if otherPoint is inside by given depth the volume
// defined by the surface in local frame.
func (p *Plane3D) isInsideLocal(otherPointLocal *Vector3D.Vector3D) bool {

	cpLocal := p.closestPointLocal(otherPointLocal)
	normalLocal := p.closestNormalLocal(otherPointLocal)
	r := otherPointLocal.Substract(cpLocal)
	return r.DotProduct(normalLocal) < 0.0
}

// NewPlane3D constructs a plane that cross \p point with surface normal \p normal.
func NewPlane3D(normal, point *Vector3D.Vector3D) *Plane3D {
	return &Plane3D{
		normal:          normal,
		point:           point,
		transform:       NewTransform3(),
		isNormalFlipped: false,
	}
}

// isBounded returns true if bounding box can be defined.
func (p *Plane3D) isBounded() bool {

	return false
}

// boundingBox returns the bounding box of this surface object.
func (p *Plane3D) boundingBox() *BoundingBox3D {

	return NewBoundingBox3D(Vector3D.NewVector(0, 0, 0), Vector3D.NewVector(0, 0, 0))
}

func (p *Plane3D) signedDistance(otherPoint *Vector3D.Vector3D) float64 {

	x := p.closestPoint(p.transform.toLocal(otherPoint))
	inside := p.isInside(otherPoint)

	sd := 0.0
	if inside {
		sd = -x.DistanceTo(otherPoint)
	} else {
		sd = x.DistanceTo(otherPoint)
	}

	if p.isNormalFlipped {
		sd = -sd
	}
	return sd
}

func (p *Plane3D) getTransform() *Transform3 {
	return p.transform
}
