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
