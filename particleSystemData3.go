package main

import "math"

// ParticleSystemData3 is the key data structure for storing particle system data. A
// single particle has position, velocity, and force attributes by default. But
// it can also have additional custom scalar or vector attributes.
type ParticleSystemData3 struct {
	radius            float64
	mass              float64
	numberOfParticles int64
	positionIdx       int64
	velocityIdx       int64
	forceIdx          int64
	scalarDataList    [][]float64
	vectorDataList    [][]*Vector3D
	neighborSearcher  *PointParallelHashGridSearcher3
	neighborLists     [][]int64
}

func NewParticleSystemData3() *ParticleSystemData3 {

	p := &ParticleSystemData3{
		radius:            0.001,
		mass:              0.001,
		numberOfParticles: 0,
		positionIdx:       0,
		velocityIdx:       0,
		forceIdx:          0,
		scalarDataList:    make([][]float64, 0),
		vectorDataList:    make([][]*Vector3D, 0),
		neighborSearcher: NewPointParallelHashGridSearcher3(
			kDefaultHashGridResolution,
			kDefaultHashGridResolution,
			kDefaultHashGridResolution,
			0.002,
		),
		neighborLists: make([][]int64, 5),
	}

	(*p).positionIdx = (*p).addVectorData()
	(*p).velocityIdx = (*p).addVectorData()
	(*p).forceIdx = (*p).addVectorData()

	return p
}

func (p *ParticleSystemData3) addVectorData() int64 {

	attrIdx := len((*p).vectorDataList)
	(*p).vectorDataList = append((*p).vectorDataList, []*Vector3D{NewVector(0, 0, 0)})
	return int64(attrIdx)

}

func (p *ParticleSystemData3) addScalarData() int64 {

	attrIdx := len((*p).scalarDataList)
	(*p).scalarDataList = append((*p).scalarDataList, []float64{0})
	return int64(attrIdx)

}

func (p *ParticleSystemData3) addParticle(newPosition, newVelocity, newForce *Vector3D) {

	newPositions := make([]*Vector3D, 0)
	newPositions = append(newPositions, newPosition)
	newVelocities := make([]*Vector3D, 0)
	newVelocities = append(newVelocities, newVelocity)
	newForces := make([]*Vector3D, 0)
	newForces = append(newForces, newForce)

	(*p).addParticles(newPositions, newVelocities, newForces)

}

func (p *ParticleSystemData3) addParticles(newPositions, newVelocities, newForces []*Vector3D) {

	var oldNumberOfParticles int64 = (*p).numberOfParticles
	var newNumberOfParticles int64 = oldNumberOfParticles + int64(len(newPositions))
	(*p).numberOfParticles = newNumberOfParticles

	p.resize(newNumberOfParticles)

	pos := (*p).positions()
	vel := (*p).velocities()
	frc := (*p).forces()

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

func (p *ParticleSystemData3) positions() []*Vector3D {

	return (*p).vectorDataList[p.positionIdx]
}

func (p *ParticleSystemData3) velocities() []*Vector3D {

	return (*p).vectorDataList[p.velocityIdx]
}

func (p *ParticleSystemData3) forces() []*Vector3D {

	return (*p).vectorDataList[p.forceIdx]
}

func (p *ParticleSystemData3) Mass() float64 {

	return (*p).mass
}

func (p *ParticleSystemData3) setMass(newMass float64) {

	(*p).mass = math.Max(newMass, 0)
}

func (p *ParticleSystemData3) resize(newNumberOfParticles int64) {

	for idx, _ := range p.scalarDataList {
		for i := int64(0); i < newNumberOfParticles-1; i++ {
			p.scalarDataList[idx] = append(p.scalarDataList[idx], 0)
		}
	}

	for idx, _ := range p.vectorDataList {

		for i := int64(0); i < newNumberOfParticles-1; i++ {
			p.vectorDataList[idx] = append(p.vectorDataList[idx], NewVector(0, 0, 0))
		}
	}
}

func (p *ParticleSystemData3) setRadius(newRadius float64) {

	p.radius = newRadius
}
