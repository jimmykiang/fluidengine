package main

import (
	"fmt"
	Vector3D "jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"log"
	"os"
	"sync"
)

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
	gravity                   *Vector3D.Vector3D
	particleSystemData        *ParticleSystemData3
	newPositions              []*Vector3D.Vector3D
	newVelocities             []*Vector3D.Vector3D
	collider                  *RigidBodyCollider3
	emitter                   *PointParticleEmitter3
	wind                      *ConstantVectorField3
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

func (p *ParticleSystemSolver3) SetEmitter(emitter *PointParticleEmitter3) {
	p.emitter = emitter
	emitter.setTarget(p.particleSystemData)
}

func NewParticleSystemSolver3() *ParticleSystemSolver3 {
	newPositions := make([]*Vector3D.Vector3D, 0)
	newPositions = append(newPositions, Vector3D.NewVector(0, 0, 0))

	newVelocities := make([]*Vector3D.Vector3D, 0)
	newVelocities = append(newVelocities, Vector3D.NewVector(0, 0, 0))

	p := &ParticleSystemSolver3{
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

// onUpdate for a standard ParticleSystemSolver3.
func (p *ParticleSystemSolver3) onUpdate(frame *Frame) {

	//numberOfFrames := frame.index - p.currentFrame.index
	//numberOfFrames := 1

	// Perform fixed time-stepping
	//for i := 0; i < numberOfFrames; i++ {

	p.beginAdvanceTimeStep(frame.timeIntervalInSeconds)

	// Add external forces.
	p.accumulateExternalForces()

	//p.timeIntegration(frame.timeIntervalInSeconds)
	p.timeIntegrationMT(frame.timeIntervalInSeconds)

	p.resolveCollision()

	p.currentFrame = frame

	// Not needed.
	//p.endAdvanceTimeStep(frame.timeIntervalInSeconds)
	//}
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

// singleThread timeIntegration
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

// mtResult collects the information from the worker threads for timeIntegrationMT.
type mtResult struct {
	newVelocity *Vector3D.Vector3D
	newPosition *Vector3D.Vector3D
}

// multiThread timeIntegration
func (p *ParticleSystemSolver3) timeIntegrationMT(timeStepsInSeconds float64) {

	n := p.particleSystemData.numberOfParticles

	threadSize := 15
	jobs := make(chan int64, n)
	results := make(chan *mtResult, n)
	var wg sync.WaitGroup

	for worker := 1; worker <= threadSize; worker++ {
		wg.Add(1)

		go func(jobs <-chan int64, results chan<- *mtResult) {
			defer wg.Done()

			for i := range jobs {
				forces := p.particleSystemData.forces()
				velocities := p.particleSystemData.velocities()
				positions := p.particleSystemData.positions()
				mass := p.particleSystemData.Mass()

				// Integrate velocity first.
				newVelocity := p.newVelocities[i]
				forceMultiply := forces[i].Multiply(timeStepsInSeconds)
				forceMultiplyDivide := forceMultiply.Divide(mass)
				newVelocity = velocities[i].Add(forceMultiplyDivide)
				//p.newVelocities[i] = newVelocity
				//p.particleSystemData.vectorDataList[p.particleSystemData.velocityIdx][i] = newVelocity

				// Integrate position.
				newPosition := p.newPositions[i]
				newVelocityMultiply := newVelocity.Multiply(timeStepsInSeconds)
				newPosition = positions[i].Add(newVelocityMultiply)
				//p.newPositions[i] = newPosition
				//p.particleSystemData.vectorDataList[p.particleSystemData.positionIdx][i] = newPosition

				results <- &mtResult{
					newVelocity: newVelocity,
					newPosition: newPosition,
				}
			}
		}(jobs, results)
	}

	for y := int64(0); y < n; y++ {

		jobs <- y
	}
	close(jobs)

	for a := int64(0); a < n; a++ {

		resultStruct := <-results

		p.newVelocities[a] = resultStruct.newVelocity
		p.particleSystemData.vectorDataList[p.particleSystemData.velocityIdx][a] = resultStruct.newVelocity

		p.newPositions[a] = resultStruct.newPosition
		p.particleSystemData.vectorDataList[p.particleSystemData.positionIdx][a] = resultStruct.newPosition
	}

	wg.Wait()
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
		relativeVel := velocities[i].Substract(p.wind.value)
		force.Add(relativeVel.Multiply(p.dragCoefficient))

		forces[i] = forces[i].Add(force)
	}
}

// beginAdvanceTimeStep is called when a time-step is about to begin.
func (p *ParticleSystemSolver3) beginAdvanceTimeStep(timeStepInSeconds float64) {

	// Clear forces.
	forces := p.particleSystemData.forces()
	for i := 0; i < len(forces); i++ {
		forces[i] = Vector3D.NewVector(0, 0, 0)
	}

	// Update collider and emitter.
	// collider does nothing.

	p.currentTime = float64(p.currentFrame.index) * p.currentFrame.timeIntervalInSeconds
	p.emitter.update(p.currentTime, timeStepInSeconds)

	// Allocate buffers.
	n := p.particleSystemData.numberOfParticles
	p.resize(n)

	p.currentTime += timeStepInSeconds
}

// initialize from the original code is only a call to onInitialize.
func (p *ParticleSystemSolver3) initialize() {

	// When initializing the solver, update the collider and emitter state as
	// well since they also affects the initial condition of the simulation.

}

func (p *ParticleSystemSolver3) setWind(wind *ConstantVectorField3) {

	p.wind = wind
}

func (p *ParticleSystemSolver3) resize(size int64) {

	for i := int64(0); i < size-1; i++ {
		p.newPositions = append(p.newPositions, Vector3D.NewVector(0, 0, 0))
		p.newVelocities = append(p.newVelocities, Vector3D.NewVector(0, 0, 0))
	}

}

func (p *ParticleSystemSolver3) saveParticleDataXyUpdate(particles *ParticleSystemData3, frame *Frame) {

	n := particles.numberOfParticles

	x := make([]float64, n)
	y := make([]float64, n)

	//positions := particles.positions()

	for i := int64(0); i < n; i++ {

		x[i] = particles.positions()[i].X
		y[i] = particles.positions()[i].Y
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	const conf = "animation/Update"
	fileNameX := fmt.Sprintf("data.#point2,%04d,x.npy", frame.index)
	fileNameY := fmt.Sprintf("data.#point2,%04d,y.npy", frame.index)

	saveNpy(path, conf, fileNameX, x, frame)
	saveNpy(path, conf, fileNameY, y, frame)
}

func (p *ParticleSystemSolver3) setIsUsingFixedSubTimeSteps(isUsing bool) {

	p.isUsingFixedSubTimeSteps = isUsing
}
