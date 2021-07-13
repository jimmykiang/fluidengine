package main

// RigidBodyCollider2 is brief 2-D rigid body collider class.
// This struct implements 2-D rigid body collider. The collider can only take
// rigid body motion with linear and rotational velocities.
type RigidBodyCollider2 struct {
	surface *Box2
}

func NewRigidBodyCollider2(surface *Box2) *RigidBodyCollider2 {
	return &RigidBodyCollider2{
		surface: surface,
	}
}

func (c *RigidBodyCollider2) update(seconds float64) {

	// do nothing?
}
