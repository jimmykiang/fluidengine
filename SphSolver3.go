package main

import (
	"fmt"
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
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
	//s.accumulateForces(timeStepInSeconds)
	//s.timeIntegration(timeStepInSeconds)
	//s.resolveCollision()
	//s.endAdvanceTimeStep(timeStepInSeconds)

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
