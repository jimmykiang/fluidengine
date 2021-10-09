package main

import "jimmykiang/fluidengine/Vector3D"

type Surface3IF interface {
	closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D
	closestDistance(point *Vector3D.Vector3D) float64
	closestNormal(point *Vector3D.Vector3D) *Vector3D.Vector3D
	getTransform() *Transform3
	isInside(position *Vector3D.Vector3D) bool
}

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
