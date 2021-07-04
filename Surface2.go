package main

type Surface2 struct {

	// Local-to-world transform.
	transform *Transform2

	// Flips normal.
	isNormalFlipped bool
}

func NewSurface2() *Surface2 {
	return &Surface2{
		transform:       NewTransform2(),
		isNormalFlipped: false,
	}
}
