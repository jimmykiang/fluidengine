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

func (s *SphSystemData2) addParticle(newPosition, newVelocity, newForce *Vector3D) {

	newPositions := make([]*Vector3D, 0)
	newPositions = append(newPositions, newPosition)
	newVelocities := make([]*Vector3D, 0)
	newVelocities = append(newVelocities, newVelocity)
	newForces := make([]*Vector3D, 0)
	newForces = append(newForces, newForce)

	(*s).addParticles(newPositions, newVelocities, newForces)

}

func (s *SphSystemData2) addParticles(newPositions, newVelocities, newForces []*Vector3D) {

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

func (s *SphSystemData2) positions() []*Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.positionIdx]
}

func (s *SphSystemData2) velocities() []*Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.velocityIdx]
}

func (s *SphSystemData2) forces() []*Vector3D {

	return (*s).particleSystemData.vectorDataList[s.particleSystemData.forceIdx]
}

func (s *SphSystemData2) buildNeighborSearcher() {
	// Use PointParallelHashGridSearcher2 by default... (now PointParallelHashGridSearcher3).
	s.particleSystemData.neighborSearcher = NewPointParallelHashGridSearcher3(
		kDefaultHashGridResolution,
		kDefaultHashGridResolution,
		0,
		2*s.kernelRadius,
	)

	size := int(s.particleSystemData.neighborSearcher.resolution.x * s.particleSystemData.neighborSearcher.resolution.y)

	for i := 0; i < size-1; i++ {
		s.particleSystemData.neighborSearcher.startIndexTable = append(
			s.particleSystemData.neighborSearcher.startIndexTable,
			int64(math.MaxInt64),
		)

		s.particleSystemData.neighborSearcher.endIndexTable = append(
			s.particleSystemData.neighborSearcher.endIndexTable,
			int64(math.MaxInt64),
		)
	}

	s.particleSystemData.neighborSearcher.build(s.positions())
}
