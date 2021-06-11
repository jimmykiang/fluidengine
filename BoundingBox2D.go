package main

import "math"

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

// Returns width of the box.
func (b *BoundingBox2D) width() float64 {

	return b.upperCorner.x - b.lowerCorner.x
}

// Returns height of the box.
func (b *BoundingBox2D) height() float64 {

	return b.upperCorner.y - b.lowerCorner.y
}
