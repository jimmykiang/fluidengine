package main

type Box2 struct {
	Surface2 *Surface2
	// Bounding box of this box.
	bound *BoundingBox2D
}

func NewBox2(boundingBox *BoundingBox2D) *Box2 {
	return &Box2{
		Surface2: NewSurface2(),
		bound:    boundingBox,
	}
}
