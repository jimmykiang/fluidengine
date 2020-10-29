package main

import "math"

// Vector3D defines simple 3-D vector data.
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
