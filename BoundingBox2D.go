package main

import (
	"math"
)

// BoundingBox2D is a 2-D axis-aligned bounding box struct.
type BoundingBox2D struct {
	// Lower corner of the bounding box.
	lowerCorner *Vector3D

	//Upper corner of the bounding box.
	upperCorner *Vector3D
}

// NewBoundingBox2D constructs a box that tightly covers two points.
func NewBoundingBox2D(point1, point2 *Vector3D) *BoundingBox2D {

	lowerCorner := NewVector(math.Min(point1.x, point2.x), math.Min(point1.y, point2.y), 0)
	upperCorner := NewVector(math.Max(point1.x, point2.x), math.Max(point1.y, point2.y), 0)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// NewBoundingBox2DReset constructs a box with the highest boundaries.
func NewBoundingBox2DReset() *BoundingBox2D {
	lowerCorner := NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	upperCorner := NewVector(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

func NewBoundingBox2DFromStruct(other *BoundingBox2D) *BoundingBox2D {

	lowerCorner := NewVector(other.lowerCorner.x, other.lowerCorner.y, other.lowerCorner.z)
	upperCorner := NewVector(other.upperCorner.x, other.upperCorner.y, other.upperCorner.z)

	return &BoundingBox2D{
		lowerCorner: lowerCorner,
		upperCorner: upperCorner,
	}
}

// Returns width of the box.
func (b *BoundingBox2D) width() float64 {

	return b.upperCorner.x - b.lowerCorner.x
}

// Returns height of the box.
func (b *BoundingBox2D) height() float64 {

	return b.upperCorner.y - b.lowerCorner.y
}

// Returns the mid-point of this box.
func (b *BoundingBox2D) midPoint() *Vector3D {

	result := b.upperCorner.Add(b.lowerCorner).Divide(2)

	return result
}

// expand this box by given delta to all direction.
// If the width of the box was x, expand(y) will result a box with
// x+y+y width.
func (b *BoundingBox2D) expand(delta float64) {

	b.lowerCorner = b.lowerCorner.Substract(NewVector(delta, delta, delta))
	b.upperCorner = b.upperCorner.Add(NewVector(delta, delta, delta))
}

// corner returns corner position. Index starts from x-first order.
func (b *BoundingBox2D) corner(idx int) *Vector3D {

	h := 0.5
	offset := make([]*Vector3D, 0, 4)

	offset = append(offset, NewVector(-h, -h, 0))
	offset = append(offset, NewVector(h, -h, 0))
	offset = append(offset, NewVector(-h, h, 0))
	offset = append(offset, NewVector(h, h, 0))

	a := NewVector(b.width(), b.height(), 0)
	//c := offset[idx].Add(b.midPoint())
	c := a.mul(offset[idx])

	result := c.Add(b.midPoint())
	return result
}

func (b *BoundingBox2D) merge(other *BoundingBox2D) {
	b.lowerCorner.x = math.Min(b.lowerCorner.x, other.lowerCorner.x)
	b.lowerCorner.y = math.Min(b.lowerCorner.y, other.lowerCorner.y)
	b.upperCorner.x = math.Max(b.upperCorner.x, other.upperCorner.x)
	b.upperCorner.y = math.Max(b.upperCorner.y, other.upperCorner.y)
}
