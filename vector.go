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

// mul computes this * (v.x, v.y, v.z).
func (v *Vector3D) mul(i *Vector3D) *Vector3D {
	return &Vector3D{
		x: v.x * i.x,
		y: v.y * i.y,
		z: v.z * i.z,
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

// Magnitude of a vector (duplicated with Length() ????)
func (v *Vector3D) LengthSquared() float64 {
	return math.Sqrt(square(v.x) + square(v.y) + square(v.z))
}

// Squared squares the components of a vector.
func (v *Vector3D) Squared() float64 {
	return square(v.x) + square(v.y) + square(v.z)
}

// Set the value of the current Vector3D with the new value from another Vector3D.
func (v *Vector3D) Set(i *Vector3D) {
	v.x = i.x
	v.y = i.y
	v.z = i.z
}

// tangential returns the tangential vector for this vector.
func (v *Vector3D) tangential() []*Vector3D {

	t := make([]*Vector3D, 0)
	var x *Vector3D
	if math.Abs(v.y) > 0 || math.Abs(v.z) > 0 {

		x = NewVector(1, 0, 0)
	} else {

		x = NewVector(0, 1, 0)
	}
	a := x.CrossProduct(v).Normalize()
	b := v.CrossProduct(a)

	t = append(t, a)
	t = append(t, b)

	return t
}

// Returns the distance to the other vector.
func (v *Vector3D) distanceTo(other *Vector3D) float64 {

	return v.Substract(other).Length()
}

func (v *Vector3D) min(o *Vector3D) *Vector3D {

	return NewVector(math.Min(v.x, o.x), math.Min(v.y, o.y), 0)
}

func (v *Vector3D) max(o *Vector3D) *Vector3D {

	return NewVector(math.Max(v.x, o.x), math.Max(v.y, o.y), 0)
}

func (v *Vector3D) isSimilar(other *Vector3D) bool {

	r := math.Abs(v.x-other.x) < kEpsilonD && math.Abs(v.y-other.y) < kEpsilonD

	return r
}
