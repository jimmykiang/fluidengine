package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"math"
)

// RigidBodyCollider2 is brief 2-D rigid body collider class.
// This struct implements 2-D rigid body collider. The collider can only take
// rigid body motion with linear and rotational velocities.
type RigidBodyCollider2 struct {
	surface             Surface2IF
	frictionCoefficient float64
	// Angular velocity of the rigid body.
	angularVelocity float64
	linearVelocity  *Vector3D.Vector3D
}

func NewRigidBodyCollider2(surface Surface2IF) *RigidBodyCollider2 {
	return &RigidBodyCollider2{
		surface:             surface,
		frictionCoefficient: 0,
		linearVelocity:      Vector3D.NewVector(0, 0, 0),
	}
}

func (c *RigidBodyCollider2) update(seconds float64) {

	// do nothing?
}

// Resolves collision for given point.
func (c *RigidBodyCollider2) resolveCollision(
	radius float64,
	restitutionCoefficient float64,
	newPosition **Vector3D.Vector3D,
	newVelocity **Vector3D.Vector3D,
) {

	colliderPoint := c.NewColliderQueryResult()

	c.getClosestPoint(c.surface, *newPosition, colliderPoint)

	// Check if the new position is penetrating the surface.
	x := c.isPenetrating(colliderPoint, *newPosition, radius)
	if x {

		// Target point is the closest non-penetrating position from the
		// new position.
		targetNormal := colliderPoint.normal
		rt := targetNormal.Multiply(radius)
		targetPoint := colliderPoint.point.Add(rt)
		colliderVelAtTargetPoint := colliderPoint.velocity

		// Get new candidate relative velocity from the target point.

		relativeVel := (*newVelocity).Substract(colliderVelAtTargetPoint)
		normalDotRelativeVel := targetNormal.DotProduct(relativeVel)
		relativeVelN := targetNormal.Multiply(normalDotRelativeVel)
		relativeVelT := relativeVel.Substract(relativeVelN)

		// Check if the velocity is facing opposite direction of the surface normal.

		if normalDotRelativeVel < 0.0 {

			// Apply restitution coefficient to the surface normal component of the velocity.

			deltaRelativeVelN := relativeVelN.Multiply(-restitutionCoefficient - 1.0)
			relativeVelN := relativeVelN.Multiply(-restitutionCoefficient)

			// Apply friction to the tangential component of the velocity From Bridson et
			// al., Robust Treatment of Collisions, Contact and Friction for Cloth Animation,
			// 2002. http://graphics.stanford.edu/papers/cloth-sig02/cloth.pdf

			if relativeVelT.LengthSquared() > 0.0 {

				a := c.frictionCoefficient * deltaRelativeVelN.Length() / relativeVelT.Length()
				frictionScale := math.Max(1-a, 0)
				relativeVelT = relativeVelT.Multiply(frictionScale)

			}
			// Reassemble the components.
			*newVelocity = relativeVelN.Add(relativeVelT).Add(colliderVelAtTargetPoint)
		}
		// Geometric fix
		//*newPosition = (*newPosition).Set(targetPoint)
		(*newPosition).Set(targetPoint)
	}
}

func (c *RigidBodyCollider2) NewColliderQueryResult() *ColliderQueryResult {
	return &ColliderQueryResult{
		distance: 0,
		point:    Vector3D.NewVector(0, 0, 0),
		normal:   Vector3D.NewVector(0, 0, 0),
		velocity: Vector3D.NewVector(0, 0, 0),
	}
}

// Outputs closest point's information.
func (c *RigidBodyCollider2) getClosestPoint(surface Surface2IF, queryPoint *Vector3D.Vector3D, result *ColliderQueryResult) {

	result.distance = surface.closestDistance(queryPoint)
	result.point = surface.closestPoint(queryPoint)
	result.normal = surface.closestNormal(queryPoint)
	result.velocity = c.velocityAt(queryPoint)
}

// Returns true if given point is in the opposite side of the surface.
func (c *RigidBodyCollider2) isPenetrating(colliderPoint *ColliderQueryResult, position *Vector3D.Vector3D, radius float64) bool {

	// If the new candidate position of the particle is inside the volume defined by
	// the surface OR the new distance to the surface is less than the particle's
	// radius, this particle is in colliding state.

	return c.surface.isInside(position) || colliderPoint.distance < radius
}

// Returns the velocity of the collider at given point.
func (c *RigidBodyCollider2) velocityAt(point *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := point.Substract(c.surface.getTransform().translation)
	a := Vector3D.NewVector(-r.Y, r.X, 0).Multiply(c.angularVelocity)
	return a.Add(c.linearVelocity)
}
