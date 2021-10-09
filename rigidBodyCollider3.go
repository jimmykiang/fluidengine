package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"math"
)

// RigidBodyCollider3 implements 3-D rigid body collider. The collider can only take
// rigid body motion with linear and rotational velocities.
type RigidBodyCollider3 struct {
	surface                  Surface3IF
	linearVelocity           *Vector3D.Vector3D
	angularVelocity          *Vector3D.Vector3D
	frictionCoefficient      float64
	onUpdateCallbackCollider OnBeginUpdateCallbackCollider
}

// ColliderQueryResult is an internal query result structure.
type ColliderQueryResult struct {
	distance float64
	point    *Vector3D.Vector3D
	normal   *Vector3D.Vector3D
	velocity *Vector3D.Vector3D
}

func (c RigidBodyCollider3) NewColliderQueryResult() *ColliderQueryResult {
	return &ColliderQueryResult{
		distance: 0,
		point:    Vector3D.NewVector(0, 0, 0),
		normal:   Vector3D.NewVector(0, 0, 0),
		velocity: Vector3D.NewVector(0, 0, 0),
	}
}

// Resolves collision for given point.
func (c RigidBodyCollider3) resolveCollision(
	radius float64,
	restitutionCoefficient float64,
	newPosition **Vector3D.Vector3D,
	newVelocity **Vector3D.Vector3D,
) {

	colliderPoint := c.NewColliderQueryResult()

	c.getClosestPoint(c.surface, *newPosition, colliderPoint)

	// Check if the new position is penetrating the surface
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

// Returns true if given point is in the opposite side of the surface.
func (c RigidBodyCollider3) isPenetrating(colliderPoint *ColliderQueryResult, position *Vector3D.Vector3D, radius float64) bool {

	// If the new candidate position of the particle is inside the volume defined by
	// the surface OR the new distance to the surface is less than the particle's
	// radius, this particle is in colliding state.

	return c.surface.isInside(position) || colliderPoint.distance < radius
}

// Outputs closest point's information.
func (c RigidBodyCollider3) getClosestPoint(surface Surface3IF, queryPoint *Vector3D.Vector3D, result *ColliderQueryResult) {

	result.distance = surface.closestDistance(queryPoint)
	result.point = surface.closestPoint(queryPoint)
	result.normal = surface.closestNormal(queryPoint)
	result.velocity = c.velocityAt(queryPoint)
}

// Returns the velocity of the collider at given point.
func (c RigidBodyCollider3) velocityAt(point *Vector3D.Vector3D) *Vector3D.Vector3D {

	//r := point.Substract(c.surface.transform.translation)
	r := point.Substract(c.surface.getTransform().translation)
	a := c.angularVelocity.CrossProduct(r)
	return c.linearVelocity.Add(a)
}

// OnBeginUpdateCallbackCollider is a brief Callback function signature type for update calls.
// This type of callback function will take the collider pointer, current
// time, and time interval in seconds.
type OnBeginUpdateCallbackCollider func(
	rigidBodyCollider *RigidBodyCollider3,
	currentTime float64,
	timeInterval float64,
)

func NewRigidBodyCollider3(surface Surface3IF) *RigidBodyCollider3 {
	return &RigidBodyCollider3{
		surface:                  surface,
		linearVelocity:           Vector3D.NewVector(0, 0, 0),
		angularVelocity:          Vector3D.NewVector(0, 0, 0),
		frictionCoefficient:      0,
		onUpdateCallbackCollider: nil,
	}
}

func (c *RigidBodyCollider3) update(seconds float64) {

	// do nothing?
}
