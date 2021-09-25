package main

import "jimmykiang/fluidengine/Vector3D"

//type ImplicitSurface3 implemented in Plane3D + Sphere3D
type ImplicitSurface3 interface {
	isBounded() bool
	boundingBox() *BoundingBox3D
	// signedDistance returns signed distance from the given point otherPoint.
	signedDistance(otherPoint *Vector3D.Vector3D) float64
	getTransform() *Transform3
}
