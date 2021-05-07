package main

import "math"

// Vector3D defines a simple 3-D vector data.
type Vector3D struct {
	x, y, z float64
}

// NewVector creates a new reference of Vector3D.
func NewVector(x, y, z float64) *Vector3D {
	return &Vector3D{x, y, z}
}

// Add Vector3D.
func (v *Vector3D) Add(i *Vector3D) *Vector3D {
	return &Vector3D{
		x: v.x + i.x,
		y: v.y + i.y,
		z: v.z + i.z,
	}
}

// Substract Vector3D.
func (v *Vector3D) Substract(i *Vector3D) *Vector3D {
	return &Vector3D{
		x: v.x - i.x,
		y: v.y - i.y,
		z: v.z - i.z,
	}
}

// Multiply Vector3D.
func (v *Vector3D) Multiply(i float64) *Vector3D {
	return &Vector3D{
		x: v.x * i,
		y: v.y * i,
		z: v.z * i,
	}
}

// Divide Vector3D.
func (v *Vector3D) Divide(i float64) *Vector3D {
	return &Vector3D{
		x: v.x / i,
		y: v.y / i,
		z: v.z / i,
	}
}

// square the value
func square(v float64) float64 {
	return math.Pow(v, 2.0)
}

// Length of a Vector3D.
func (v *Vector3D) Length() float64 {
	return math.Sqrt(square(v.x) + square(v.y) + square(v.z))
}

// Normalize a Vector3D.
func (v *Vector3D) Normalize() *Vector3D {
	length := v.Length()
	if length == 0.0 {
		return v
	}
	return NewVector(v.x/length, v.y/length, v.z/length)
}

// DotProduct from 2 Vector3D.
func (v *Vector3D) DotProduct(o *Vector3D) float64 {
	return ((v.x * o.x) + (v.y * o.y) + (v.z * o.z))
}

// CrossProduct from 2 vectors (tuple with w == 0).
func (v *Vector3D) CrossProduct(o *Vector3D) *Vector3D {
	return NewVector((*v).y*o.z-v.z*o.y, v.z*o.x-v.x*o.z, v.x*o.y-v.y*o.x)
}

// Magnitude of a vector
func (v *Vector3D) LengthSquared() float64 {
	return math.Sqrt(square(v.x) + square(v.y) + square(v.z))
}

// Set the value of the current Vector3D with the new value from another Vector3D.
func (v *Vector3D) Set(i *Vector3D) {
	v.x = i.x
	v.y = i.y
	v.z = i.z
}
