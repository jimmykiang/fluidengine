package main

// ConstantVectorField3 3-D constant vector field.
type ConstantVectorField3 struct {

	value *Vector3D
}

func NewConstantVectorField3() *ConstantVectorField3 {
	return &ConstantVectorField3{

		value: NewVector(0,0,0),
	}
}

func (c *ConstantVectorField3) withValue(v *Vector3D) {

	c.value.Set(v)
}

