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
