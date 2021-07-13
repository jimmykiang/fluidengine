package main

type Transform2 struct {
	translation *Vector3D
	orientation float64
	cosAngle    float64
	sinAngle    float64
}

func NewTransform2() *Transform2 {
	return &Transform2{
		translation: NewVector(0, 0, 0),
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

		a := bboxInWorld.lowerCorner.min(cornerInWorld)
		bboxInWorld.lowerCorner = NewVector(a.x, a.y, 0)

		b := bboxInWorld.upperCorner.max(cornerInWorld)
		bboxInWorld.upperCorner = NewVector(b.x, b.y, 0)
	}

	return bboxInWorld
}

// Transforms a point in local space to the world coordinate.
func (t *Transform2) toWorldPointInLocal(pointInLocal *Vector3D) *Vector3D {

	// Convert to the world frame.
	x := t.cosAngle*pointInLocal.x - t.sinAngle*pointInLocal.y + t.translation.x
	y := t.sinAngle*pointInLocal.x + t.cosAngle*pointInLocal.y + t.translation.y

	return NewVector(x, y, 0)
}

// Transforms a point in world coordinate to the local frame.
func (t *Transform2) toLocal(pointInWorld *Vector3D) *Vector3D {

	// Convert to the local frame.
	xmt := pointInWorld.Substract(t.translation)

	x := t.cosAngle*xmt.x + t.sinAngle*xmt.y
	y := -t.sinAngle*xmt.x + t.cosAngle*xmt.y
	return NewVector(x, y, 0)
}
