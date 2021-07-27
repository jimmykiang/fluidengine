package main

import "jimmykiang/fluidengine/Vector3D"

// ConstantVectorField3 3-D constant vector field.
type ConstantVectorField3 struct {
	value *Vector3D.Vector3D
}

func NewConstantVectorField3() *ConstantVectorField3 {
	return &ConstantVectorField3{

		value: Vector3D.NewVector(0, 0, 0),
	}
}

func (c *ConstantVectorField3) withValue(v *Vector3D.Vector3D) {

	c.value.Set(v)
}
