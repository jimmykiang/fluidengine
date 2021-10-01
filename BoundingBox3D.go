package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"math"
)

// BoundingBox3D is a 2-D axis-aligned bounding box struct.
type BoundingBox3D struct {
	// Lower corner of the bounding box.
	lowerCorner *Vector3D.Vector3D

	//Upper corner of the bounding box.
	upperCorner *Vector3D.Vector3D
}

// NewBoundingBox3D constructs a box that tightly covers two points.
func NewBoundingBox3D(point1, point2 *Vector3D.Vector3D) *BoundingBox3D {

	lowerCorner := Vector3D.NewVector(math.Min(point1.X, point2.X), math.Min(point1.Y, point2.Y), math.Min(point1.Z, point2.Z))
	upperCorner := Vector3D.NewVector(math.Max(point1.X, point2.X), math.Max(point1.Y, point2.Y), math.Max(point1.Z, point2.Z))

	return &BoundingBox3D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// Returns width of the box.
func (b *BoundingBox3D) width() float64 {

	return b.upperCorner.X - b.lowerCorner.X
}

// Returns height of the box.
func (b *BoundingBox3D) height() float64 {

	return b.upperCorner.Y - b.lowerCorner.Y
}

// Returns depth of the box.
func (b *BoundingBox3D) depth() float64 {

	return b.upperCorner.Z - b.lowerCorner.Z
}

// NewBoundingBox3DReset constructs a box with the highest boundaries.
func NewBoundingBox3DReset() *BoundingBox3D {
	lowerCorner := Vector3D.NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	upperCorner := Vector3D.NewVector(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)

	return &BoundingBox3D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// corner returns corner position. Index starts from x-first order.
func (b *BoundingBox3D) corner(idx int) *Vector3D.Vector3D {

	h := 0.5
	offset := make([]*Vector3D.Vector3D, 0, 8)

	offset = append(offset, Vector3D.NewVector(-h, -h, -h))
	offset = append(offset, Vector3D.NewVector(h, -h, -h))
	offset = append(offset, Vector3D.NewVector(-h, h, -h))
	offset = append(offset, Vector3D.NewVector(h, h, -h))
	offset = append(offset, Vector3D.NewVector(-h, -h, h))
	offset = append(offset, Vector3D.NewVector(h, -h, h))
	offset = append(offset, Vector3D.NewVector(-h, h, h))
	offset = append(offset, Vector3D.NewVector(h, h, h))

	a := Vector3D.NewVector(b.width(), b.height(), b.depth())
	//c := offset[idx].Add(b.midPoint())
	c := a.Mul(offset[idx])

	result := c.Add(b.midPoint())
	return result
}

// Returns the mid-point of this box.
func (b *BoundingBox3D) midPoint() *Vector3D.Vector3D {

	result := b.upperCorner.Add(b.lowerCorner).Divide(2)

	return result
}

// NewBoundingBox3DFromStruct constructs a box with other box instance.
func NewBoundingBox3DFromStruct(other *BoundingBox3D) *BoundingBox3D {

	lowerCorner := Vector3D.NewVector(other.lowerCorner.X, other.lowerCorner.Y, other.lowerCorner.Z)
	upperCorner := Vector3D.NewVector(other.upperCorner.X, other.upperCorner.Y, other.upperCorner.Z)

	return &BoundingBox3D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// expand this box by given delta to all direction.
// If the width of the box was x, expand(y) will result a box with
// x+y+y width.
func (b *BoundingBox3D) expand(delta float64) {

	b.lowerCorner = b.lowerCorner.Substract(Vector3D.NewVector(delta, delta, delta))
	b.upperCorner = b.upperCorner.Add(Vector3D.NewVector(delta, delta, delta))
}

func (b *BoundingBox3D) contains(point *Vector3D.Vector3D) bool {

	if b.upperCorner.X < point.X || b.lowerCorner.X > point.X {
		return false
	}
	if b.upperCorner.Y < point.Y || b.lowerCorner.Y > point.Y {
		return false
	}
	if b.upperCorner.Z < point.Z || b.lowerCorner.Z > point.Z {
		return false
	}
	return true
}
