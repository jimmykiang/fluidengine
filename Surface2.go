package main

import "jimmykiang/fluidengine/Vector3D"

type Surface2IF interface {
	closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D
	closestDistance(point *Vector3D.Vector3D) float64
	closestNormal(point *Vector3D.Vector3D) *Vector3D.Vector3D
	getTransform() *Transform2
	isInside(position *Vector3D.Vector3D) bool
}

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
