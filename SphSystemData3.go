package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"math"
)

// SphSystemData3 is a 3-D SPH particle system data.
// Extends ParticleSystemData3 to specialize the data model for SPH.
// It includes density and pressure array as a default particle attribute, and
// it also contains SPH utilities such as interpolation operator.
type SphSystemData3 struct {
	particleSystemData *ParticleSystemData3

	// Target density of this particle system in kg/m^3.
	targetDensity float64
	// Target spacing of this particle system in meters.
	targetSpacing float64
	// Relative radius of SPH kernel.
	// SPH kernel radius divided by target spacing.0
	kernelRadiusOverTargetSpacing float64
	//SPH kernel radius in meters.
	kernelRadius float64
	pressureIdx  int64
	densityIdx   int64
}

func NewSphSystemData3() *SphSystemData3 {
	s := &SphSystemData3{
		particleSystemData:            NewParticleSystemData3(),
		targetDensity:                 constants.KWaterDensity,
		targetSpacing:                 0.2,
		kernelRadiusOverTargetSpacing: 1.8,
		kernelRadius:                  1,
		pressureIdx:                   0,
		densityIdx:                    0,
	}

	s.densityIdx = (*s).particleSystemData.addScalarData()
	s.pressureIdx = (*s).particleSystemData.addScalarData()
	s.setTargetSpacing(s.targetSpacing)

	return s
}

func (s *SphSystemData3) setTargetDensity(targetDensity float64) {

	s.targetDensity = targetDensity
	s.computeMass()
}

func (s *SphSystemData3) setTargetSpacing(spacing float64) {

	s.particleSystemData.setRadius(spacing)
	s.targetSpacing = spacing
	s.kernelRadius = s.kernelRadiusOverTargetSpacing * spacing
	s.computeMass()
}

func (s *SphSystemData3) computeMass() {

	points := make([]*Vector3D.Vector3D, 0)
	pointsGenerator := NewBccLatticePointGenerator()

	sampleBound := NewBoundingBox3D(
		Vector3D.NewVector(-1.5*s.kernelRadius, -1.5*s.kernelRadius, -1.5*s.kernelRadius),
		Vector3D.NewVector(1.5*s.kernelRadius, 1.5*s.kernelRadius, 1.5*s.kernelRadius),
	)
	pointsGenerator.generate(sampleBound, s.targetSpacing, &points)

	maxNumberDensity := 0.0
	kernel := NewSphStdKernel3(s.kernelRadius)

	for i := 0; i < len(points); i++ {
		point := points[i]
		sum := 0.0

		for j := 0; j < len(points); j++ {

			neighborPoint := points[j]
			x := neighborPoint.DistanceTo(point)
			sum += kernel.operatorKernel(x)
		}
		maxNumberDensity = math.Max(maxNumberDensity, sum)
	}

	newMass := s.targetDensity / maxNumberDensity
	s.particleSystemData.setMass(newMass)
}

func (s *SphSystemData3) addParticles(newPositions, newVelocities, newForces []*Vector3D.Vector3D) {

	var oldNumberOfParticles int64 = (*s).particleSystemData.numberOfParticles
	var newNumberOfParticles int64 = oldNumberOfParticles + int64(len(newPositions))
	(*s).particleSystemData.numberOfParticles = newNumberOfParticles

	s.particleSystemData.resize(newNumberOfParticles)

	pos := (*s).positions()
	vel := (*s).velocities()
	frc := (*s).forces()

	if (len(newPositions)) > 0 {
		for i := 0; i < len(newPositions); i++ {

			pos[int64(i)+oldNumberOfParticles] = newPositions[i]
		}
	}

	if (len(newVelocities)) > 0 {
		for i := 0; i < len(newPositions); i++ {

			vel[int64(i)+oldNumberOfParticles] = newVelocities[i]
		}
	}

	if (len(newPositions)) > 0 {
		for i := 0; i < len(newForces); i++ {

			frc[int64(i)+oldNumberOfParticles] = newForces[i]
		}
	}
}

func (s *SphSystemData3) positions() []*Vector3D.Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.positionIdx]
}

func (s *SphSystemData3) velocities() []*Vector3D.Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.velocityIdx]
}

func (s *SphSystemData3) forces() []*Vector3D.Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.forceIdx]
}
