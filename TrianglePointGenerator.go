package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"math"
)

type TrianglePointGenerator struct {
}

func NewTrianglePointGenerator() *TrianglePointGenerator {
	return &TrianglePointGenerator{}
}

func (t *TrianglePointGenerator) generate(
	boundingBox *BoundingBox2D,
	spacing float64,
	points *([]*Vector3D.Vector3D),
) {
	t.forEachPoint(
		boundingBox,
		spacing,
		points,
		t.callback,
	)

}

func (t *TrianglePointGenerator) callback(points *([]*Vector3D.Vector3D), v *Vector3D.Vector3D) bool {
	*points = append(*points, Vector3D.NewVector(v.X, v.Y, v.Z))
	return true
}

func (t *TrianglePointGenerator) forEachPoint(
	boundingBox *BoundingBox2D,
	spacing float64,
	points *[]*Vector3D.Vector3D,
	callback func(*([]*Vector3D.Vector3D), *Vector3D.Vector3D) bool,
) {

	halfSpacing := spacing / 2
	ySpacing := spacing * math.Sqrt(3) / 2
	boxWidth := boundingBox.width()
	boxHeight := boundingBox.height()

	position := Vector3D.NewVector(0, 0, 0)
	hasOffset := false
	shouldQuit := false

	for j := float64(0); j*ySpacing <= boxHeight && !shouldQuit; j++ {

		position.Y = j*ySpacing + boundingBox.lowerCorner.Y
		var offset float64
		if hasOffset {

			offset = halfSpacing
		} else {
			offset = 0
		}

		for i := float64(0); i*spacing+offset <= boxWidth && !shouldQuit; i++ {
			position.X = i*spacing + offset + boundingBox.lowerCorner.X
			if !callback(points, position) {
				shouldQuit = true
				break
			}
		}
		hasOffset = !hasOffset
	}
}
