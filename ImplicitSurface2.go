package main

type ImplicitSurface2 struct {

	// Local-to-world transform.
	transform       *Transform2
	isNormalFlipped bool
}

func NewImplicitSurface2() *ImplicitSurface2 {
	return &ImplicitSurface2{
		transform:       NewTransform2(),
		isNormalFlipped: false,
	}
}
