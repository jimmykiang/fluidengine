package Vector3D

import (
	"jimmykiang/fluidengine/constants"
	"math"
)

// Vector3D defines a simple 3-D vector data.
type Vector3D struct {
	X, Y, Z float64
}

// NewVector creates a new reference of Vector3D.
func NewVector(x, y, z float64) *Vector3D {
	return &Vector3D{x, y, z}
}

// Add Vector3D.
func (v *Vector3D) Add(i *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X + i.X,
		Y: v.Y + i.Y,
		Z: v.Z + i.Z,
	}
}

// Substract Vector3D.
func (v *Vector3D) Substract(i *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X - i.X,
		Y: v.Y - i.Y,
		Z: v.Z - i.Z,
	}
}

// Multiply Vector3D.
func (v *Vector3D) Multiply(i float64) *Vector3D {
	return &Vector3D{
		X: v.X * i,
		Y: v.Y * i,
		Z: v.Z * i,
	}
}

// Mul computes this * (v.x, v.y, v.z).
func (v *Vector3D) Mul(i *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X * i.X,
		Y: v.Y * i.Y,
		Z: v.Z * i.Z,
	}
}

// Divide Vector3D.
func (v *Vector3D) Divide(i float64) *Vector3D {
	return &Vector3D{
		X: v.X / i,
		Y: v.Y / i,
		Z: v.Z / i,
	}
}

// square the value
func square(v float64) float64 {
	return math.Pow(v, 2.0)
}

// Length of a Vector3D.
func (v *Vector3D) Length() float64 {
	return math.Sqrt(square(v.X) + square(v.Y) + square(v.Z))
}

// Normalize a Vector3D.
func (v *Vector3D) Normalize() *Vector3D {
	length := v.Length()
	if length == 0.0 {
		return v
	}
	return NewVector(v.X/length, v.Y/length, v.Z/length)
}

// DotProduct from 2 Vector3D.
func (v *Vector3D) DotProduct(o *Vector3D) float64 {
	return ((v.X * o.X) + (v.Y * o.Y) + (v.Z * o.Z))
}

// CrossProduct from 2 vectors (tuple with w == 0).
func (v *Vector3D) CrossProduct(o *Vector3D) *Vector3D {
	return NewVector((*v).Y*o.Z-v.Z*o.Y, v.Z*o.X-v.X*o.Z, v.X*o.Y-v.Y*o.X)
}

// Magnitude of a vector (duplicated with Length() ????)
func (v *Vector3D) LengthSquared() float64 {
	return math.Sqrt(square(v.X) + square(v.Y) + square(v.Z))
}

// Squared squares the components of a vector.
func (v *Vector3D) Squared() float64 {
	return square(v.X) + square(v.Y) + square(v.Z)
}

// Set the value of the current Vector3D with the new value from another Vector3D.
func (v *Vector3D) Set(i *Vector3D) {
	v.X = i.X
	v.Y = i.Y
	v.Z = i.Z
}

// tangential returns the tangential vector for this vector.
func (v *Vector3D) Tangential() []*Vector3D {

	t := make([]*Vector3D, 0)
	var x *Vector3D
	if math.Abs(v.Y) > 0 || math.Abs(v.Z) > 0 {

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
func (v *Vector3D) DistanceTo(other *Vector3D) float64 {

	return v.Substract(other).Length()
}

func (v *Vector3D) Min(o *Vector3D) *Vector3D {

	return NewVector(math.Min(v.X, o.X), math.Min(v.Y, o.Y), 0)
}

func (v *Vector3D) Max(o *Vector3D) *Vector3D {

	return NewVector(math.Max(v.X, o.X), math.Max(v.Y, o.Y), 0)
}

func (v *Vector3D) IsSimilar(other *Vector3D) bool {

	r := math.Abs(v.X-other.X) < constants.KEpsilonD && math.Abs(v.Y-other.Y) < constants.KEpsilonD

	return r
}

func (v *Vector3D) DistanceSquaredTo(other *Vector3D) float64 {

	return v.Substract(other).Squared()
}
