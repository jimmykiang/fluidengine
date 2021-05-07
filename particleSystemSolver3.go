package main

// ParticleSystemSolver3 is a Basic 3-D particle system solver.
//
// This class implements basic particle system solver. It includes gravity,
// air drag, and collision. But it does not compute particle-to-particle
// interaction. Thus, this solver is suitable for performing simple spray-like
// simulations with low computational cost. This class can be further extend
// to add more sophisticated simulations, such as SPH, to handle
// particle-to-particle intersection.
type ParticleSystemSolver3 struct {
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
	emitter                   *PointParticleEmitter3
	wind                      *Vector3D
}

func (p *ParticleSystemSolver3) ParticleSystemData() *ParticleSystemData3 {
	return p.particleSystemData
}

func (p *ParticleSystemSolver3) DragCoefficient() float64 {
	return p.dragCoefficient
}

func (p *ParticleSystemSolver3) SetDragCoefficient(dragCoefficient float64) {
	p.dragCoefficient = dragCoefficient
}

func (p *ParticleSystemSolver3) RestitutionCoefficient() float64 {
	return p.restitutionCoefficient
}

func (p *ParticleSystemSolver3) SetRestitutionCoefficient(restitutionCoefficient float64) {
	p.restitutionCoefficient = restitutionCoefficient
}

func (p *ParticleSystemSolver3) Collider() *RigidBodyCollider3 {
	return p.collider
}

func (p *ParticleSystemSolver3) SetCollider(collider *RigidBodyCollider3) {
	p.collider = collider
}

func NewParticleSystemSolver3() *ParticleSystemSolver3 {
	newPositions := make([]*Vector3D, 0)
	newPositions = append(newPositions, NewVector(0, 0, 0))

	newVelocities := make([]*Vector3D, 0)
	newVelocities = append(newVelocities, NewVector(0, 0, 0))

	p := &ParticleSystemSolver3{
		currentFrame:              NewFrame(),
		isUsingFixedSubTimeSteps:  true,
		numberOfFixedSubTimeSteps: 1,
		currentTime:               0,
		dragCoefficient:           0,
		restitutionCoefficient:    0,
		gravity:                   NewVector(0, kGravity, 0),
		particleSystemData:        NewParticleSystemData3(),
		newPositions:              newPositions,
		newVelocities:             newVelocities,
		collider:                  nil,
		emitter:                   nil,
		wind:                      NewVector(0, 0, 0),
	}

	p.currentFrame.index = -1
	return p
}

// onUpdate for a standard ParticleSystemSolver3.
func (p *ParticleSystemSolver3) onUpdate(frame *Frame) {

	//numberOfFrames := frame.index - p.currentFrame.index
	numberOfFrames := 1

	// Perform fixed time-stepping
	for i := 0; i < numberOfFrames; i++ {

		p.beginAdvanceTimeStep(frame.timeIntervalInSeconds)

		// Add external forces.
		p.accumulateExternalForces()

		p.timeIntegration(frame.timeIntervalInSeconds)

		p.resolveCollision()

		// Not needed.
		//p.endAdvanceTimeStep(frame.timeIntervalInSeconds)
	}
	//p.currentFrame = frame
}

func (p *ParticleSystemSolver3) endAdvanceTimeStep(seconds float64) {
	// Update data.

	// Not needed.
}

func (p *ParticleSystemSolver3) resolveCollision() {

	numberOfParticles := p.particleSystemData.numberOfParticles
	radius := p.particleSystemData.radius

	for i := 0; i < int(numberOfParticles); i++ {
		p.collider.resolveCollision(radius, p.restitutionCoefficient, &p.newPositions[i], &p.newVelocities[i])
		p.particleSystemData.vectorDataList[p.particleSystemData.velocityIdx][i] = p.newVelocities[i]
		p.particleSystemData.vectorDataList[p.particleSystemData.positionIdx][i] = p.newPositions[i]
	}
}

func (p *ParticleSystemSolver3) timeIntegration(timeStepsInSeconds float64) {

	n := p.particleSystemData.numberOfParticles
	forces := p.particleSystemData.forces()
	velocities := p.particleSystemData.velocities()
	positions := p.particleSystemData.positions()
	mass := p.particleSystemData.Mass()

	for i := 0; i < int(n); i++ {

		// Integrate velocity first.
		newVelocity := p.newVelocities[i]
		forceMultiply := forces[i].Multiply(timeStepsInSeconds)
		forceMultiplyDivide := forceMultiply.Divide(mass)
		newVelocity = velocities[i].Add(forceMultiplyDivide)
		p.newVelocities[i] = newVelocity
		p.particleSystemData.vectorDataList[p.particleSystemData.velocityIdx][i] = newVelocity
		//(*p.particleSystemData.vectorDataList[p.particleSystemData.velocityIdx][i]).Set(newVelocity)

		// Integrate position.

		newPosition := p.newPositions[i]
		newVelocityMultiply := newVelocity.Multiply(timeStepsInSeconds)
		newPosition = positions[i].Add(newVelocityMultiply)
		p.newPositions[i] = newPosition
		p.particleSystemData.vectorDataList[p.particleSystemData.positionIdx][i] = newPosition
		//(*p.particleSystemData.vectorDataList[p.particleSystemData.positionIdx][i]).Set(newPosition)
	}
}

func (p *ParticleSystemSolver3) accumulateExternalForces() {

	n := p.particleSystemData.numberOfParticles
	forces := p.particleSystemData.forces()
	velocities := p.particleSystemData.velocities()
	mass := p.particleSystemData.Mass()

	for i := 0; i < int(n); i++ {
		// Gravity.
		force := p.gravity.Multiply(mass)

		// Wind forces.
		relativeVel := velocities[i]
		force.Add(relativeVel.Multiply(p.dragCoefficient))

		forces[i] = forces[i].Add(force)
	}
}

// beginAdvanceTimeStep is called when a time-step is about to begin.
func (p *ParticleSystemSolver3) beginAdvanceTimeStep(timeStepInSeconds float64) {

	// Clear forces.
	forces := p.particleSystemData.forces()
	for i := 0; i < len(forces); i++ {
		forces[i] = NewVector(0, 0, 0)
	}

	// Update collider and emitter.
	// Do nothing

	// Allocate buffers.
	// Do nothing

	_ = forces
}

// initialize from the original code is only a call to onInitialize.
func (p *ParticleSystemSolver3) initialize() {

	// When initializing the solver, update the collider and emitter state as
	// well since they also affects the initial condition of the simulation.

}