package main

type Surface3 struct {

	// Local-to-world transform.
	transform *Transform3

	// Flips normal.
	isNormalFlipped bool
}

func NewSurface3() *Surface3 {
	return &Surface3{
		transform:       NewTransform3(),
		isNormalFlipped: false,
	}
}
