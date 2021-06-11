package main

import "math"

// SphSystemData2 is a 2-D SPH particle system data.
// Extends ParticleSystemData2 to specialize the data model for SPH.
// It includes density and pressure array as a default particle attribute, and
// it also contains SPH utilities such as interpolation operator.
type SphSystemData2 struct {
	particleSystemData *ParticleSystemData3

	// Target density of this particle system in kg/m^2.
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

func NewSphSystemData2() *SphSystemData2 {
	s := &SphSystemData2{
		particleSystemData:            NewParticleSystemData3(),
		targetDensity:                 kWaterDensity,
		targetSpacing:                 0.1,
		kernelRadiusOverTargetSpacing: 1.8,
		kernelRadius:                  0,
		pressureIdx:                   0,
		densityIdx:                    0,
	}

	s.densityIdx = (*s).particleSystemData.addScalarData()
	s.pressureIdx = (*s).particleSystemData.addScalarData()
	s.setTargetSpacing(s.targetSpacing)

	return s
}

func (s *SphSystemData2) setTargetSpacing(spacing float64) {

	s.particleSystemData.setRadius(spacing)
	s.targetSpacing = spacing
	s.kernelRadius = s.kernelRadiusOverTargetSpacing * spacing
	s.computeMass()
}

func (s *SphSystemData2) computeMass() {

	points := make([]*Vector3D, 0)
	pointsGenerator := NewTrianglePointGenerator()

	sampleBound := NewBoundingBox2D(
		NewVector(-1.5*s.kernelRadius, -1.5*s.kernelRadius, 0),
		NewVector(1.5*s.kernelRadius, 1.5*s.kernelRadius, 0),
	)
	pointsGenerator.generate(sampleBound, s.targetSpacing, &points)

	maxNumberDensity := 0.0
	kernel := NewSphStdKernel2(s.kernelRadius)

	for i := 0; i < len(points); i++ {
		point := points[i]
		sum := 0.0

		for j := 0; j < len(points); j++ {

			neighborPoint := points[j]
			x := neighborPoint.distanceTo(point)
			sum += kernel.operatorKernel(x)
		}
		maxNumberDensity = math.Max(maxNumberDensity, sum)
	}

	newMass := s.targetDensity / maxNumberDensity
	s.particleSystemData.setMass(newMass)
}
func (s *SphSystemData2) setTargetDensity(targetDensity float64) {

	s.targetDensity = targetDensity
	s.computeMass()
}
