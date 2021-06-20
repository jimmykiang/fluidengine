package main

// ParticleSystemSolver2 is a basic 2-D particle system solver.
// This struct implements basic particle system solver. It includes gravity,
// air drag, and collision. But it does not compute particle-to-particle
// interaction. Thus, this solver is suitable for performing simple spray-like
// simulations with low computational cost. This class can be further extend
// to add more sophisticated simulations, such as SPH, to handle
// particle-to-particle intersection.
// SphSolver2
type ParticleSystemSolver2 struct {
	currentFrame              *Frame
	isUsingFixedSubTimeSteps  bool
	numberOfFixedSubTimeSteps int64
	currentTime               float64
	dragCoefficient           float64
	restitutionCoefficient    float64
	gravity                   *Vector3D
	particleSystemData        *ParticleSystemData3
	newPositions              []*Vector3D
	newVelocities             []*Vector3D
	collider                  *RigidBodyCollider3
	emitter                   *VolumeParticleEmitter2
	wind                      *ConstantVectorField3
}

func NewParticleSystemSolver2() *ParticleSystemSolver2 {
	newPositions := make([]*Vector3D, 0)
	newPositions = append(newPositions, NewVector(0, 0, 0))

	newVelocities := make([]*Vector3D, 0)
	newVelocities = append(newVelocities, NewVector(0, 0, 0))

	p := &ParticleSystemSolver2{
		currentFrame:              NewFrame(),
		isUsingFixedSubTimeSteps:  true,
		numberOfFixedSubTimeSteps: 1,
		currentTime:               0,
		dragCoefficient:           0.0001,
		restitutionCoefficient:    0,
		gravity:                   NewVector(0, kGravity, 0),
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

func (p *ParticleSystemSolver2) setIsUsingFixedSubTimeSteps(isUsing bool) {

	p.isUsingFixedSubTimeSteps = isUsing
}
