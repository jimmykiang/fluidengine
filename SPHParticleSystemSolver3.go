package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
)

// SPHParticleSystemSolver3 is a Basic 3-D particle system solver.
//
// This class implements basic particle system solver. It includes gravity,
// air drag, and collision. But it does not compute particle-to-particle
// interaction. Thus, this solver is suitable for performing simple spray-like
// simulations with low computational cost. This class can be further extend
// to add more sophisticated simulations, such as SPH, to handle
// particle-to-particle intersection.
// Used by SphSolver3
type SPHParticleSystemSolver3 struct {
	currentFrame              *Frame
	isUsingFixedSubTimeSteps  bool
	numberOfFixedSubTimeSteps int64
	currentTime               float64
	dragCoefficient           float64
	restitutionCoefficient    float64
	gravity                   *Vector3D.Vector3D
	particleSystemData        *ParticleSystemData3
	newPositions              []*Vector3D.Vector3D
	newVelocities             []*Vector3D.Vector3D
	collider                  *RigidBodyCollider3
	emitter                   *VolumeParticleEmitter3
	wind                      *ConstantVectorField3
}

func NewSPHParticleSystemSolver3() *SPHParticleSystemSolver3 {
	newPositions := make([]*Vector3D.Vector3D, 0)
	newPositions = append(newPositions, Vector3D.NewVector(0, 0, 0))

	newVelocities := make([]*Vector3D.Vector3D, 0)
	newVelocities = append(newVelocities, Vector3D.NewVector(0, 0, 0))

	p := &SPHParticleSystemSolver3{
		currentFrame:              NewFrame(),
		isUsingFixedSubTimeSteps:  true,
		numberOfFixedSubTimeSteps: 1,
		currentTime:               0,
		dragCoefficient:           0,
		restitutionCoefficient:    0,
		gravity:                   Vector3D.NewVector(0, constants.KGravity, 0),
		particleSystemData:        NewParticleSystemData3(),
		newPositions:              newPositions,
		newVelocities:             newVelocities,
		collider:                  nil,
		emitter:                   nil,
		wind:                      NewConstantVectorField3(),
	}

	p.currentFrame.index = -1
	return p
}

func (p *SPHParticleSystemSolver3) setIsUsingFixedSubTimeSteps(isUsing bool) {

	p.isUsingFixedSubTimeSteps = isUsing
}
