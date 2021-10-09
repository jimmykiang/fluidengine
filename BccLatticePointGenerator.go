package main

import (
	"jimmykiang/fluidengine/Vector3D"
)

// BccLatticePointGenerator is a Body-centered lattice points generator.
// http://en.wikipedia.org/wiki/Cubic_crystal_system
// http://mathworld.wolfram.com/CubicClosePacking.html
type BccLatticePointGenerator struct {
}

func NewBccLatticePointGenerator() *BccLatticePointGenerator {
	return &BccLatticePointGenerator{}
}

func (b *BccLatticePointGenerator) generate(
	boundingBox *BoundingBox3D,
	spacing float64,
	points *([]*Vector3D.Vector3D),
) {
	b.forEachPoint(
		boundingBox,
		spacing,
		points,
		b.callback,
	)

}

func (b *BccLatticePointGenerator) callback(points *([]*Vector3D.Vector3D), v *Vector3D.Vector3D) bool {
	*points = append(*points, Vector3D.NewVector(v.X, v.Y, v.Z))
	return true
}

// forEachPoint iterates every BCC-lattice points inside \p boundingBox
// where \p spacing is the size of the unit cell of BCC structure.
func (b *BccLatticePointGenerator) forEachPoint(
	boundingBox *BoundingBox3D,
	spacing float64,
	points *[]*Vector3D.Vector3D,
	callback func(*([]*Vector3D.Vector3D), *Vector3D.Vector3D) bool,
) {

	halfSpacing := spacing / 2
	boxWidth := boundingBox.width()
	boxHeight := boundingBox.height()
	boxDepth := boundingBox.depth()

	position := Vector3D.NewVector(0, 0, 0)
	hasOffset := false
	shouldQuit := false

	for k := float64(0); k*halfSpacing <= boxDepth && !shouldQuit; k++ {

		position.Z = k*halfSpacing + boundingBox.lowerCorner.Z
		var offset float64
		if hasOffset {

			offset = halfSpacing
		} else {
			offset = 0
		}

		for j := float64(0); j*spacing+offset <= boxHeight && !shouldQuit; j++ {
			position.Y = j*spacing + offset + boundingBox.lowerCorner.Y

			for i := float64(0); i*spacing+offset <= boxWidth; i++ {
				position.X = i*spacing + offset + boundingBox.lowerCorner.X

				if !callback(points, position) {
					shouldQuit = true
					break
				}
			}
		}
		hasOffset = !hasOffset
	}
}
