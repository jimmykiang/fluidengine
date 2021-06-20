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

func NewBoundingBox2DFromStruct(other *BoundingBox2D) *BoundingBox2D{

	lowerCorner:=NewVector(other.lowerCorner.x, other.lowerCorner.y, other.lowerCorner.z)
	upperCorner := NewVector(other.upperCorner.x, other.upperCorner.y, other.upperCorner.z)

	return &BoundingBox2D{
		lowerCorner : lowerCorner,
		upperCorner : upperCorner,
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

	return b.upperCorner.Add(b.lowerCorner).Divide(2)
}

// expand this box by given delta to all direction.
// If the width of the box was x, expand(y) will result a box with
// x+y+y width.
func (b *BoundingBox2D) expand(delta float64) {

	b.lowerCorner = b.lowerCorner.Substract(NewVector(delta,delta,delta))
	b.upperCorner = b.upperCorner.Add(NewVector(delta,delta,delta))
}
