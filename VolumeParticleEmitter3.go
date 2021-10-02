package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"math"
	"math/rand"
)

// VolumeParticleEmitter3 is a 3-D volumetric particle emitter.
// emits particles from volumetric geometry.
// Constructs an emitter that spawns particles from given implicit surface
// which defines the volumetric geometry. Provided bounding box limits
// the particle generation region.
type VolumeParticleEmitter3 struct {
	implicitSurface          *ImplicitSurfaceSet3
	particles                *SphSystemData3
	bounds                   *BoundingBox3D
	spacing                  float64
	initialVel               *Vector3D.Vector3D
	linearVel                *Vector3D.Vector3D
	angularVel               *Vector3D.Vector3D
	maxNumberOfParticles     float64
	jitter                   float64
	isOneShot                bool
	allowOverlapping         bool
	seed                     int64
	isEnabled                bool
	pointsGen                *BccLatticePointGenerator
	numberOfEmittedParticles float64
}

func NewVolumeParticleEmitter3(
	implicitSurface *ImplicitSurfaceSet3,
	maxRegion *BoundingBox3D,
	spacing float64,
	initialVel *Vector3D.Vector3D,
) *VolumeParticleEmitter3 {
	return &VolumeParticleEmitter3{
		implicitSurface:          implicitSurface,
		bounds:                   maxRegion,
		spacing:                  spacing,
		initialVel:               initialVel,
		linearVel:                Vector3D.NewVector(0, 0, 0),
		angularVel:               Vector3D.NewVector(0, 0, 0),
		maxNumberOfParticles:     constants.KMaxSize,
		numberOfEmittedParticles: 0,
		jitter:                   0,
		isOneShot:                true,
		allowOverlapping:         false,
		seed:                     0,
		isEnabled:                true,
		pointsGen:                NewBccLatticePointGenerator(),
	}
}

func (e *VolumeParticleEmitter3) setTarget(particles *SphSystemData3) {

	e.particles = particles
	e.onSetTarget(particles)
}

func (e *VolumeParticleEmitter3) onSetTarget(particles *SphSystemData3) {

	// Do nothing.
}

func (e *VolumeParticleEmitter3) onUpdate() {

	particles := e.particles

	if !e.isEnabled {
		return
	}

	newPositions := make([]*Vector3D.Vector3D, 0, 0)
	newVelocities := make([]*Vector3D.Vector3D, 0, 0)

	e.emit(particles, &newPositions, &newVelocities)

	particles.addParticles(newPositions, newVelocities, nil)

	if e.isOneShot {
		e.isEnabled = false
	}
}

func (e *VolumeParticleEmitter3) emit(particles *SphSystemData3, newPositions, newVelocities *[]*Vector3D.Vector3D) {

	e.implicitSurface.updateQueryEngine()

	region := NewBoundingBox3D(e.bounds.lowerCorner, e.bounds.upperCorner)

	if e.implicitSurface.isBounded() {
		//todo: surfaceBBox:=
	}

	// Reserving more space for jittering
	j := e.jitter
	maxJitterDist := 0.5 * j * e.spacing
	numNewParticles := 0.0

	callback := func(points *([]*Vector3D.Vector3D), point *Vector3D.Vector3D) bool {

		randomDir := e.uniformSampleSphere(rand.Float64(), rand.Float64())
		offset := randomDir.Multiply(maxJitterDist)
		candidate := point.Add(offset)

		if e.implicitSurface.signedDistance(candidate) <= 0.0 {
			if e.numberOfEmittedParticles < e.maxNumberOfParticles {
				*newPositions = append(*newPositions, candidate)
				e.numberOfEmittedParticles++
				numNewParticles++
			} else {
				return false
			}
		}
		return true
	}

	_ = callback

	if e.allowOverlapping || e.isOneShot {

		e.pointsGen.forEachPoint(region, e.spacing, nil, callback)
	} else {
		// Use serial hash grid searcher for continuous update.
		// todo.
	}

	///////////////////////////
	// not needed?
	//*newVelocities = make([]*Vector3D.Vector3D, len(*newPositions), len(*newPositions))

	// original code from: \FluidEngine\fluid-engine-dev\src\jet\volume_particle_emitter2.cpp
	//newVelocities->parallelForEachIndex([&](size_t i) {
	//	(*newVelocities)[i] = velocityAt((*newPositions)[i]);
	//});

	/////////////////////////////

	e.parallelForEachIndex(newVelocities, newPositions)
}

func (e *VolumeParticleEmitter3) uniformSampleSphere(u1 float64, u2 float64) *Vector3D.Vector3D {
	y := 1 - 2*u1
	r := math.Sqrt(math.Max(0, 1-y*y))
	phi := math.Pi * 2 * u2
	x := r * math.Cos(phi)
	z := r * math.Sin(phi)

	return Vector3D.NewVector(x, y, z)
}

func (e *VolumeParticleEmitter3) parallelForEachIndex(newVelocities, newPositions *[]*Vector3D.Vector3D) {
	for i := 0.0; i < e.numberOfEmittedParticles; i++ {

		e.callback(i, newVelocities, newPositions)
	}
}

func (e *VolumeParticleEmitter3) callback(i float64, newVelocities, newPositions *[]*Vector3D.Vector3D) {

	*newVelocities = append(*newVelocities, e.velocityAt((*newPositions)[int64(i)]))
}

func (e *VolumeParticleEmitter3) velocityAt(point *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := point.Substract(e.implicitSurface.surfaces[0].getTransform().translation)
	a := e.angularVel.CrossProduct(r)

	return a.Add(e.linearVel).Add(e.initialVel)
}
