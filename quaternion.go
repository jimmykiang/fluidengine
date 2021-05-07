package main

// Quaternion struct defined as q = w + xi + yj + zk.
type Quaternion struct {
	// Real part.
	w float32

	// Imaginary parts.
	x, y, z float32
}

// newQuaternion creates an identity Quaternion.
func newQuaternion() *Quaternion {
	return &Quaternion{1, 0, 0, 0}
}
