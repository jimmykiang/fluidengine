package main

import (
	Vector3D "jimmykiang/fluidengine/Vector3D"
	"math"
)

// BoundingBox2D is a 2-D axis-aligned bounding box struct.
type BoundingBox2D struct {
	// Lower corner of the bounding box.
	lowerCorner *Vector3D.Vector3D

	//Upper corner of the bounding box.
	upperCorner *Vector3D.Vector3D
}

// NewBoundingBox2D constructs a box that tightly covers two points.
func NewBoundingBox2D(point1, point2 *Vector3D.Vector3D) *BoundingBox2D {

	lowerCorner := Vector3D.NewVector(math.Min(point1.X, point2.X), math.Min(point1.Y, point2.Y), 0)
	upperCorner := Vector3D.NewVector(math.Max(point1.X, point2.X), math.Max(point1.Y, point2.Y), 0)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// NewBoundingBox2DReset constructs a box with the highest boundaries.
func NewBoundingBox2DReset() *BoundingBox2D {
	lowerCorner := Vector3D.NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	upperCorner := Vector3D.NewVector(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

func NewBoundingBox2DFromStruct(other *BoundingBox2D) *BoundingBox2D {

	lowerCorner := Vector3D.NewVector(other.lowerCorner.X, other.lowerCorner.Y, other.lowerCorner.Z)
	upperCorner := Vector3D.NewVector(other.upperCorner.X, other.upperCorner.Y, other.upperCorner.Z)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// Returns width of the box.
func (b *BoundingBox2D) width() float64 {

	return b.upperCorner.X - b.lowerCorner.X
}

// Returns height of the box.
func (b *BoundingBox2D) height() float64 {

	return b.upperCorner.Y - b.lowerCorner.Y
}

// Returns the mid-point of this box.
func (b *BoundingBox2D) midPoint() *Vector3D.Vector3D {

	result := b.upperCorner.Add(b.lowerCorner).Divide(2)

	return result
}

// expand this box by given delta to all direction.
// If the width of the box was x, expand(y) will result a box with
// x+y+y width.
func (b *BoundingBox2D) expand(delta float64) {

	b.lowerCorner = b.lowerCorner.Substract(Vector3D.NewVector(delta, delta, delta))
	b.upperCorner = b.upperCorner.Add(Vector3D.NewVector(delta, delta, delta))
}

// corner returns corner position. Index starts from x-first order.
func (b *BoundingBox2D) corner(idx int) *Vector3D.Vector3D {

	h := 0.5
	offset := make([]*Vector3D.Vector3D, 0, 4)

	offset = append(offset, Vector3D.NewVector(-h, -h, 0))
	offset = append(offset, Vector3D.NewVector(h, -h, 0))
	offset = append(offset, Vector3D.NewVector(-h, h, 0))
	offset = append(offset, Vector3D.NewVector(h, h, 0))

	a := Vector3D.NewVector(b.width(), b.height(), 0)
	//c := offset[idx].Add(b.midPoint())
	c := a.Mul(offset[idx])

	result := c.Add(b.midPoint())
	return result
}

func (b *BoundingBox2D) merge(other *BoundingBox2D) {
	b.lowerCorner.X = math.Min(b.lowerCorner.X, other.lowerCorner.X)
	b.lowerCorner.Y = math.Min(b.lowerCorner.Y, other.lowerCorner.Y)
	b.upperCorner.X = math.Max(b.upperCorner.X, other.upperCorner.X)
	b.upperCorner.Y = math.Max(b.upperCorner.Y, other.upperCorner.Y)
}

func (b *BoundingBox2D) contains(point *Vector3D.Vector3D) bool {

	if b.upperCorner.X < point.X || b.lowerCorner.X > point.X {
		return false
	}
	if b.upperCorner.Y < point.Y || b.lowerCorner.Y > point.Y {
		return false
	}
	return true
}
