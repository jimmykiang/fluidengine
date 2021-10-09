package main

import "jimmykiang/fluidengine/Vector3D"

//type ImplicitSurface2 implemented in Plane2D + Sphere2D
type ImplicitSurface2 interface {
	isBounded() bool
	boundingBox() *BoundingBox2D
	signedDistance(otherPoint *Vector3D.Vector3D) float64
	getTransform() *Transform2
}
