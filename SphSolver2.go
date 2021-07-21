package main

import "math"

// SphSolver2 implements a 2-D SPH solver. The main pressure solver is based on
// equation-of-state (EOS).
type SphSolver2 struct {
	particleSystemData    *SphSystemData2
	particleSystemSolver2 *ParticleSystemSolver2
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
}

func NewSphSolver2() *SphSolver2 {
	s := &SphSolver2{
		particleSystemSolver2:      NewParticleSystemSolver2(),
		particleSystemData:         NewSphSystemData2(),
		wind:                       NewConstantVectorField3(),
		eosExponent:                7.0,
		negativePressureScale:      0,
		viscosityCoefficient:       0.01,
		pseudoViscosityCoefficient: 10,
		speedOfSound:               100,
		timeStepLimitScale:         1,
	}

	s.particleSystemSolver2.setIsUsingFixedSubTimeSteps(false)
	return s
}

func (s *SphSolver2) setPseudoViscosityCoefficient(newPseudoViscosityCoefficient float64) {

	s.pseudoViscosityCoefficient = math.Max(newPseudoViscosityCoefficient, 0)
}

func (s *SphSolver2) setEmitter(newEmitter *VolumeParticleEmitter2) {

	s.particleSystemSolver2.emitter = newEmitter
	newEmitter.setTarget(s.particleSystemData)
}

func (s *SphSolver2) setCollider(collider *RigidBodyCollider2) {

	s.particleSystemSolver2.SetCollider(collider)
}

func (s *SphSolver2) onUpdate(frame *Frame) {

	s.onInitialize()

	s.advanceTimeStep(frame.timeIntervalInSeconds)

}

// onInitialize initializes the simulator.
func (s *SphSolver2) onInitialize() {
	// When initializing the solver, update the collider and emitter state as
	// well since they also affects the initial condition of the simulation.

	s.updateEmitter(0.0)
}

func (s *SphSolver2) numberOfSubTimeSteps(timeIntervalInSeconds float64) int64 {

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

	timeStepLimitBySpeed := kTimeStepLimitBySpeedFactor * kernelRadius / s.speedOfSound
	timeStepLimitByForce := kTimeStepLimitByForceFactor * math.Sqrt(kernelRadius*mass/maxForceMagnitude)
	desiredTimeStep := s.timeStepLimitScale * math.Min(timeStepLimitBySpeed, timeStepLimitByForce)
	return int64(math.Ceil(timeIntervalInSeconds / desiredTimeStep))
}

func (s *SphSolver2) updateEmitter(f float64) {

	s.particleSystemSolver2.emitter.onUpdate()
}

func (s *SphSolver2) advanceTimeStep(timeIntervalInSeconds float64) {

	// Perform adaptive time-stepping
	remainingTime := timeIntervalInSeconds

	if remainingTime > kEpsilonD {

		numSteps := s.numberOfSubTimeSteps(remainingTime)
		actualTimeInterval := remainingTime / float64(numSteps)

		s.onAdvanceTimeStep(actualTimeInterval)

	}
}

func (s *SphSolver2) onAdvanceTimeStep(timeStepInSeconds float64) {

	s.beginAdvanceTimeStep(timeStepInSeconds)

	s.accumulateForces(timeStepInSeconds)

}

func (s *SphSolver2) updateCollider(timeStepInSeconds float64) {
	s.particleSystemSolver2.collider.update(timeStepInSeconds)
}

func (s *SphSolver2) resize(size int64) {

	for i := int64(0); i < size-1; i++ {
		s.particleSystemSolver2.newPositions = append(s.particleSystemSolver2.newPositions, NewVector(0, 0, 0))
		s.particleSystemSolver2.newVelocities = append(s.particleSystemSolver2.newVelocities, NewVector(0, 0, 0))
	}
}

func (s *SphSolver2) beginAdvanceTimeStep(timeStepInSeconds float64) {

	// Clear forces.
	forces := s.particleSystemData.forces()
	for i := 0; i < len(forces); i++ {
		forces[i] = NewVector(0, 0, 0)
	}

	// Update collider and emitter.
	s.updateCollider(timeStepInSeconds)
	s.particleSystemSolver2.emitter.onUpdate()

	// Allocate buffers.
	n := s.particleSystemData.particleSystemData.numberOfParticles
	s.resize(n)

	s.onBeginAdvanceTimeStep(timeStepInSeconds)
}

func (s *SphSolver2) onBeginAdvanceTimeStep(seconds float64) {

	particles := s.particleSystemData
	particles.buildNeighborSearcher()
	particles.buildNeighborLists()
	particles.updateDensities()
}

func (s *SphSolver2) accumulateForces(timeStepInSeconds float64) {

	s.accumulateNonPressureForces(timeStepInSeconds)
}

func (s *SphSolver2) accumulateNonPressureForces(timeStepInSeconds float64) {

	s.accumulateExternalForces(timeStepInSeconds)
	s.accumulateViscosityForce()
}

func (s *SphSolver2) accumulateExternalForces(timeStepInSeconds float64) {

	n := s.particleSystemData.particleSystemData.numberOfParticles
	forces := s.particleSystemData.particleSystemData.forces()
	velocities := s.particleSystemData.particleSystemData.velocities()
	mass := s.particleSystemData.particleSystemData.Mass()

	for i := 0; i < int(n); i++ {
		// Gravity.
		force := s.particleSystemSolver2.gravity.Multiply(mass)

		// Wind forces.
		relativeVel := velocities[i].Substract(s.particleSystemSolver2.wind.value)
		force.Add(relativeVel.Multiply(s.particleSystemSolver2.dragCoefficient))

		forces[i] = forces[i].Add(force)
	}
}

func (s *SphSolver2) accumulateViscosityForce() {
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
			dist := x[i].distanceTo(x[j])

			a := s.viscosityCoefficient * massSquared * kernel.secondDerivative(dist)
			b := v[j].Substract(v[i])
			c := b.Divide(d[j])
			f[i] = f[i].Add(c.Multiply(a))
		}
	}
}
