package main

import "jimmykiang/fluidengine/Vector3D"

// Transform3 represents 3-D rigid body transform.
type Transform3 struct {
	translation            *Vector3D.Vector3D
	orientation            *Quaternion
	orientationMat3        Matrix
	inverseOrientationMat3 Matrix
}

// Transforms a point in world coordinate to the local frame.
func (t Transform3) toLocal(pointInWorld *Vector3D.Vector3D) *Vector3D.Vector3D {

	a := pointInWorld.Substract(t.translation)
	return t.inverseOrientationMat3.MultiplyMatrixByTuple(a)
}

// Transforms a point in local space to the world coordinate.
func (t Transform3) toWorld(pointInLocal *Vector3D.Vector3D) *Vector3D.Vector3D {

	a := t.inverseOrientationMat3.MultiplyMatrixByTuple(pointInLocal)
	return a.Add(t.translation)
}

// Transforms a direction in local space to the world coordinate.
func (t Transform3) toWorldDirection(dirInLocal *Vector3D.Vector3D) *Vector3D.Vector3D {

	return t.orientationMat3.MultiplyMatrixByTuple(dirInLocal)
}

func NewTransform3() *Transform3 {
	return &Transform3{
		translation:            Vector3D.NewVector(0, 0, 0),
		orientation:            newQuaternion(),
		orientationMat3:        New3x3IdentityMatrix(),
		inverseOrientationMat3: New3x3IdentityMatrix(),
	}
}
