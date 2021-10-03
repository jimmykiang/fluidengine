package main

import (
	"fmt"
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"jimmykiang/fluidengine/mathHelper"
	"jimmykiang/fluidengine/physicsHelper"
	"log"
	"math"
	"os"
)

// SphSolver3 implements a 3-D SPH solver. The main pressure solver is based on
// equation-of-state (EOS).
type SphSolver3 struct {
	particleSystemData    *SphSystemData3
	particleSystemSolver3 *SPHParticleSystemSolver3
	wind                  *ConstantVectorField3
	// Exponent component of equation-of-state (or Tait's equation).
	eosExponent float64
	// Negative pressure scaling factor. Zero means clamping. One means do nothing.
	negativePressureScale float64
	// Viscosity coefficient.
	viscosityCoefficient float64
	//Pseudo-viscosity coefficient velocity filtering.
	// This is a minimum "safety-net" for SPH solver which is quite sensitive to the parameters.
	pseudoViscosityCoefficient float64
	// Speed of sound in medium to determine the stiffness of the system.
	// Ideally, it should be the actual speed of sound in the fluid, but in
	// practice, use lower value to trace-off performance and compressibility.
	speedOfSound float64
	// Scales the max allowed time-step.
	timeStepLimitScale float64
	currentFrame       *Frame
}

func NewSphSolver3() *SphSolver3 {
	s := &SphSolver3{
		particleSystemSolver3:      NewSPHParticleSystemSolver3(),
		particleSystemData:         NewSphSystemData3(),
		wind:                       NewConstantVectorField3(),
		eosExponent:                7.0,
		negativePressureScale:      0,
		viscosityCoefficient:       0.01,
		pseudoViscosityCoefficient: 10,
		speedOfSound:               100,
		timeStepLimitScale:         1,
		currentFrame:               NewFrame(),
	}

	s.particleSystemSolver3.setIsUsingFixedSubTimeSteps(false)
	s.currentFrame.index = -1
	return s
}

func (s *SphSolver3) setPseudoViscosityCoefficient(newPseudoViscosityCoefficient float64) {

	s.pseudoViscosityCoefficient = math.Max(newPseudoViscosityCoefficient, 0)
}
func (s *SphSolver3) setEmitter(newEmitter *VolumeParticleEmitter3) {

	s.particleSystemSolver3.emitter = newEmitter
	newEmitter.setTarget(s.particleSystemData)
}

func (s *SphSolver3) setCollider(collider *RigidBodyCollider3) {

	s.particleSystemSolver3.SetCollider(collider)
}

func (s *SphSolver3) setViscosityCoefficient(f float64) {

	s.viscosityCoefficient = f
}

func (s *SphSolver3) onUpdate(frame *Frame) {
	if s.currentFrame.index < 0 {
		s.onInitialize()
	}

	s.advanceTimeStep(frame.timeIntervalInSeconds)
	s.currentFrame = frame
}

// onInitialize initializes the simulator.
func (s *SphSolver3) onInitialize() {
	// When initializing the solver, update the collider and emitter state as
	// well since they also affects the initial condition of the simulation.

	s.updateEmitter(0.0)
}

func (s *SphSolver3) updateEmitter(f float64) {

	s.particleSystemSolver3.emitter.onUpdate()
}

func (p *SphSolver3) saveParticleDataXyUpdate(particles *ParticleSystemData3, frame *Frame) {

	n := particles.numberOfParticles

	x := make([]float64, n)
	y := make([]float64, n)

	for i := int64(0); i < n; i++ {

		x[i] = particles.positions()[i].X
		y[i] = particles.positions()[i].Y
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	const conf = "animation/SphSolver3WaterDrop"
	fileNameX := fmt.Sprintf("data.#point2,%04d,x.npy", frame.index)
	fileNameY := fmt.Sprintf("data.#point2,%04d,y.npy", frame.index)

	saveNpy(path, conf, fileNameX, x, frame)
	saveNpy(path, conf, fileNameY, y, frame)
}

func (s *SphSolver3) advanceTimeStep(timeIntervalInSeconds float64) {

	// Perform adaptive time-stepping
	remainingTime := timeIntervalInSeconds

	for remainingTime > constants.KEpsilonD {

		numSteps := s.numberOfSubTimeSteps(remainingTime)
		actualTimeInterval := remainingTime / float64(numSteps)

		//println("numSteps:", numSteps)
		s.onAdvanceTimeStep(actualTimeInterval)
		remainingTime -= actualTimeInterval
	}
}

func (s *SphSolver3) numberOfSubTimeSteps(timeIntervalInSeconds float64) int64 {

	//particles := NewSphSystemData2()
	particles := s.particleSystemData
	numberOfParticles := particles.particleSystemData.numberOfParticles
	f := particles.forces()

	kernelRadius := particles.kernelRadius
	mass := particles.particleSystemData.mass
	maxForceMagnitude := 0.0

	for i := int64(0); i < numberOfParticles; i++ {
		maxForceMagnitude = math.Max(maxForceMagnitude, f[i].Length())
	}

	timeStepLimitBySpeed := constants.KTimeStepLimitBySpeedFactor * kernelRadius / s.speedOfSound
	timeStepLimitByForce := constants.KTimeStepLimitByForceFactor * math.Sqrt(kernelRadius*mass/maxForceMagnitude)
	desiredTimeStep := s.timeStepLimitScale * math.Min(timeStepLimitBySpeed, timeStepLimitByForce)
	return int64(math.Ceil(timeIntervalInSeconds / desiredTimeStep))
}

func (s *SphSolver3) onAdvanceTimeStep(timeStepInSeconds float64) {

	s.beginAdvanceTimeStep(timeStepInSeconds)
	s.accumulateForces(timeStepInSeconds)
	s.timeIntegration(timeStepInSeconds)
	s.resolveCollision()
	s.endAdvanceTimeStep(timeStepInSeconds)
}

func (s *SphSolver3) timeIntegration(timeStepsInSeconds float64) {

	n := s.particleSystemData.particleSystemData.numberOfParticles
	forces := s.particleSystemData.forces()
	velocities := s.particleSystemData.velocities()
	positions := s.particleSystemData.positions()
	mass := s.particleSystemData.particleSystemData.Mass()

	for i := 0; i < int(n); i++ {

		// Integrate velocity first.
		newVelocity := s.particleSystemSolver3.newVelocities[i]
		forceMultiply := forces[i].Multiply(timeStepsInSeconds)
		forceMultiplyDivide := forceMultiply.Divide(mass)
		newVelocity = velocities[i].Add(forceMultiplyDivide)
		s.particleSystemSolver3.newVelocities[i] = newVelocity

		// Integrate position.
		newPosition := s.particleSystemSolver3.newPositions[i]
		newVelocityMultiply := newVelocity.Multiply(timeStepsInSeconds)
		newPosition = positions[i].Add(newVelocityMultiply)
		s.particleSystemSolver3.newPositions[i] = newPosition
	}
}

func (s *SphSolver3) beginAdvanceTimeStep(timeStepInSeconds float64) {

	// Clear forces.
	forces := s.particleSystemData.forces()
	for i := 0; i < len(forces); i++ {
		forces[i] = Vector3D.NewVector(0, 0, 0)
	}

	// Update collider and emitter.
	s.updateCollider(timeStepInSeconds)
	s.particleSystemSolver3.emitter.onUpdate()

	// Allocate buffers.
	n := s.particleSystemData.particleSystemData.numberOfParticles
	s.resize(n)
	s.particleSystemSolver3.particleSystemData.resize(n)

	s.onBeginAdvanceTimeStep(timeStepInSeconds)
}

func (s *SphSolver3) updateCollider(timeStepInSeconds float64) {
	s.particleSystemSolver3.collider.update(timeStepInSeconds)
}

func (s *SphSolver3) resize(size int64) {

	for i := int64(0); i < size-1; i++ {
		s.particleSystemSolver3.newPositions = append(s.particleSystemSolver3.newPositions, Vector3D.NewVector(0, 0, 0))
		s.particleSystemSolver3.newVelocities = append(s.particleSystemSolver3.newVelocities, Vector3D.NewVector(0, 0, 0))
	}
}

func (s *SphSolver3) onBeginAdvanceTimeStep(seconds float64) {
	particles := s.particleSystemData
	particles.buildNeighborSearcher()
	particles.buildNeighborLists()
	particles.updateDensities()
}

func (s *SphSolver3) accumulateForces(timeStepInSeconds float64) {

	s.accumulateNonPressureForces(timeStepInSeconds)
	s.accumulatePressureForce(timeStepInSeconds)
}

func (s *SphSolver3) accumulateNonPressureForces(timeStepInSeconds float64) {

	s.accumulateExternalForces(timeStepInSeconds)
	s.accumulateViscosityForce()
}

func (s *SphSolver3) accumulateViscosityForce() {
	particles := s.particleSystemData.particleSystemData
	numberOfParticles := s.particleSystemData.particleSystemData.numberOfParticles
	x := s.particleSystemData.positions()
	v := s.particleSystemData.velocities()
	d := s.particleSystemData.densities()
	f := s.particleSystemData.forces()

	massSquared := math.Pow(s.particleSystemData.particleSystemData.mass, 2)

	kernel := NewSphSpikyKernel2(s.particleSystemData.kernelRadius)

	for i := int64(0); i < numberOfParticles; i++ {

		neighbors := particles.neighborLists[i]

		for _, j := range neighbors {
			dist := x[i].DistanceTo(x[j])

			a := s.viscosityCoefficient * massSquared * kernel.secondDerivative(dist)
			b := v[j].Substract(v[i])
			c := b.Divide(d[j])
			f[i] = f[i].Add(c.Multiply(a))
		}
	}
}

func (s *SphSolver3) accumulatePressureForce(timeStepInSeconds float64) {

	x := s.particleSystemData.positions()
	d := s.particleSystemData.densities()
	p := s.particleSystemData.pressures()
	f := s.particleSystemData.forces()

	s.computePressure()
	s.accumulatePressureForceInternal(x, d, p, f)
}

func (s *SphSolver3) accumulatePressureForceInternal(
	positions []*Vector3D.Vector3D,
	densities []float64,
	pressures []float64,
	pressureForces []*Vector3D.Vector3D,
) {
	particles := s.particleSystemData.particleSystemData
	numberOfParticles := particles.numberOfParticles
	massSquared := particles.Mass() * particles.Mass()
	kernel := NewSphSpikyKernel3(s.particleSystemData.kernelRadius)

	for i := int64(0); i < numberOfParticles; i++ {
		neighbors := particles.neighborLists[i]
		for _, j := range neighbors {
			dist := positions[i].DistanceTo(positions[j])

			if dist > 0.0 {
				a := positions[j].Substract(positions[i])
				dir := a.Divide(dist)
				b := massSquared * (pressures[i]/(densities[i]*densities[i]) +
					pressures[j]/(densities[j]*densities[j]))
				c := kernel.gradient(dist, dir)
				d := c.Multiply(b)
				pressureForces[i] = pressureForces[i].Substract(d)
			}
		}
	}
}

func (s *SphSolver3) computePressure() {
	particles := s.particleSystemData.particleSystemData
	numberOfParticles := particles.numberOfParticles
	d := s.particleSystemData.densities()
	p := s.particleSystemData.pressures()

	// See Murnaghan-Tait equation of state from
	// https://en.wikipedia.org/wiki/Tait_equation
	targetDensity := s.particleSystemData.targetDensity
	eosScale := targetDensity * s.speedOfSound * s.speedOfSound

	for i := int64(0); i < numberOfParticles; i++ {
		p[i] = physicsHelper.ComputePressureFromEos(
			d[i],
			targetDensity,
			eosScale,
			s.eosExponent,
			s.negativePressureScale,
		)
	}
}

func (s *SphSolver3) accumulateExternalForces(timeStepInSeconds float64) {

	n := s.particleSystemData.particleSystemData.numberOfParticles
	forces := s.particleSystemData.particleSystemData.forces()
	velocities := s.particleSystemData.particleSystemData.velocities()
	mass := s.particleSystemData.particleSystemData.Mass()

	for i := 0; i < int(n); i++ {
		// Gravity.
		force := s.particleSystemSolver3.gravity.Multiply(mass)

		// Wind forces.
		relativeVel := velocities[i].Substract(s.particleSystemSolver3.wind.value)
		//force.Add(relativeVel.Multiply(-s.particleSystemSolver2.dragCoefficient))
		force = force.Add(relativeVel.Multiply(-s.particleSystemSolver3.dragCoefficient))

		forces[i] = forces[i].Add(force)
	}
}

func (s *SphSolver3) resolveCollision() {

	numberOfParticles := s.particleSystemData.particleSystemData.numberOfParticles
	radius := s.particleSystemData.particleSystemData.radius

	for i := 0; i < int(numberOfParticles); i++ {
		s.particleSystemSolver3.collider.resolveCollision(
			radius,
			s.particleSystemSolver3.restitutionCoefficient,
			&s.particleSystemSolver3.newPositions[i],
			&s.particleSystemSolver3.newVelocities[i],
		)
	}
}

func (s *SphSolver3) endAdvanceTimeStep(timeStepInSeconds float64) {
	// Update data.
	n := s.particleSystemData.particleSystemData.numberOfParticles
	positions := s.particleSystemData.positions()
	velocities := s.particleSystemData.velocities()

	for i := 0; i < int(n); i++ {

		positions[i] = s.particleSystemSolver3.newPositions[i]
		velocities[i] = s.particleSystemSolver3.newVelocities[i]
	}

	s.onEndAdvanceTimeStep(timeStepInSeconds)
}

func (s *SphSolver3) onEndAdvanceTimeStep(timeStepInSeconds float64) {
	s.computePseudoViscosity(timeStepInSeconds)
	numberOfParticles := s.particleSystemData.particleSystemData.numberOfParticles
	densities := s.particleSystemData.densities()

	maxDensity := 0.0

	for i := 0; i < int(numberOfParticles); i++ {
		maxDensity = math.Max(maxDensity, densities[i])
	}
}

func (s *SphSolver3) computePseudoViscosity(timeStepInSeconds float64) {

	particles := s.particleSystemData
	//particles := s.particleSystemData.particleSystemData
	numberOfParticles := s.particleSystemData.particleSystemData.numberOfParticles
	x := particles.positions()
	d := particles.densities()
	v := particles.velocities()
	mass := particles.particleSystemData.mass
	kernel := NewSphSpikyKernel3(s.particleSystemData.kernelRadius)

	smoothedVelocities := make([]*Vector3D.Vector3D, 0, 0)

	for i := 0; i < int(numberOfParticles); i++ {

		weightSum := 0.0
		smoothedVelocity := Vector3D.NewVector(0, 0, 0)
		neighbors := s.particleSystemData.particleSystemData.neighborLists[i]

		for _, j := range neighbors {
			dist := x[i].DistanceTo(x[j])
			wj := mass / d[j] * kernel.operatorKernel(dist)
			weightSum += wj

			a := v[j].Multiply(wj)
			smoothedVelocity = smoothedVelocity.Add(a)
		}

		wi := mass / d[i]
		weightSum += wi
		a := v[i].Multiply(wi)
		smoothedVelocity = smoothedVelocity.Add(a)

		if weightSum > 0.0 {
			smoothedVelocity = smoothedVelocity.Divide(weightSum)
		}
		smoothedVelocities = append(smoothedVelocities, smoothedVelocity)
	}

	factor := timeStepInSeconds * s.pseudoViscosityCoefficient
	factor = mathHelper.Clamp(
		factor,
		0,
		1,
	)

	for i := int64(0); i < numberOfParticles; i++ {

		v[i] = mathHelper.Lerp(v[i], smoothedVelocities[i], factor)
	}
}
