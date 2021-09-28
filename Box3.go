package main

type Box3 struct {
	Surface3 *Surface3
	// Bounding box of this box.
	bound *BoundingBox3D
	// Local-to-world transform.
	transform       *Transform3
	isNormalFlipped bool
}

func NewBox3(boundingBox *BoundingBox3D) *Box3 {
	return &Box3{
		Surface3:        NewSurface3(),
		bound:           boundingBox,
		transform:       NewTransform3(),
		isNormalFlipped: true,
	}
}
