package main

import "math"

// SphSolver3 implements a 3-D SPH solver. The main pressure solver is based on
// equation-of-state (EOS).
type SphSolver3 struct {
	particleSystemData    *SphSystemData3
	particleSystemSolver3 *ParticleSystemSolver3
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
		particleSystemSolver3:      NewParticleSystemSolver3(),
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

//func (s *SphSolver3) setEmitter(newEmitter *VolumeParticleEmitter3) {
//
//	s.particleSystemSolver3.emitter = newEmitter
//	newEmitter.setTarget(s.particleSystemData)
//}
