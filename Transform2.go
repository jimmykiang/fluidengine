package main

import Vector3D "jimmykiang/fluidengine/Vector3D"

type Transform2 struct {
	translation *Vector3D.Vector3D
	orientation float64
	cosAngle    float64
	sinAngle    float64
}

func NewTransform2() *Transform2 {
	return &Transform2{
		translation: Vector3D.NewVector(0, 0, 0),
		orientation: 0,
		cosAngle:    1,
		sinAngle:    0,
	}
}

// toWorld transforms a bounding box in local space to the world coordinate.
func (t *Transform2) toWorld(bboxInLocal *BoundingBox2D) *BoundingBox2D {

	bboxInWorld := NewBoundingBox2DReset()

	for i := 0; i < 4; i++ {
		cornerInWorld := t.toWorldPointInLocal(bboxInLocal.corner(i))

		a := bboxInWorld.lowerCorner.Min(cornerInWorld)
		bboxInWorld.lowerCorner = Vector3D.NewVector(a.X, a.Y, 0)

		b := bboxInWorld.upperCorner.Max(cornerInWorld)
		bboxInWorld.upperCorner = Vector3D.NewVector(b.X, b.Y, 0)
	}

	return bboxInWorld
}

// toWorld transforms a bounding box in local space to the world coordinate.
func (t *Transform2) toWorldArgVector(pointInLocal *Vector3D.Vector3D) *Vector3D.Vector3D {

	return Vector3D.NewVector(
		(t.cosAngle*pointInLocal.X)-(t.sinAngle*pointInLocal.Y+t.translation.X),
		(t.sinAngle*pointInLocal.X)-(t.cosAngle*pointInLocal.Y+t.translation.Y),
		0,
	)
}

// Transforms a point in local space to the world coordinate.
func (t *Transform2) toWorldPointInLocal(pointInLocal *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Convert to the world frame.
	x := t.cosAngle*pointInLocal.X - t.sinAngle*pointInLocal.Y + t.translation.X
	y := t.sinAngle*pointInLocal.X + t.cosAngle*pointInLocal.Y + t.translation.Y

	return Vector3D.NewVector(x, y, 0)
}

// Transforms a point in world coordinate to the local frame.
func (t *Transform2) toLocal(pointInWorld *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Convert to the local frame.
	xmt := pointInWorld.Substract(t.translation)

	x := t.cosAngle*xmt.X + t.sinAngle*xmt.Y
	y := -t.sinAngle*xmt.X + t.cosAngle*xmt.Y
	return Vector3D.NewVector(x, y, 0)
}

func (t *Transform2) toWorldDirection(dirInLocal *Vector3D.Vector3D) *Vector3D.Vector3D {
	// Convert to the world frame.

	x := t.cosAngle*dirInLocal.X - t.sinAngle*dirInLocal.Y
	y := t.sinAngle*dirInLocal.X + t.cosAngle*dirInLocal.Y
	return Vector3D.NewVector(x, y, 0)
}
