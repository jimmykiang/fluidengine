package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"math"
	"math/rand"
)

// VolumeParticleEmitter2 is a 2-D volumetric particle emitter.
// emits particles from volumetric geometry.
// Constructs an emitter that spawns particles from given implicit surface
// which defines the volumetric geometry. Provided bounding box limits
// the particle generation region.
type VolumeParticleEmitter2 struct {
	implicitSurface          *ImplicitSurfaceSet2
	particles                *SphSystemData2
	bounds                   *BoundingBox2D
	spacing                  float64
	initialVel               *Vector3D.Vector3D
	linearVel                *Vector3D.Vector3D
	angularVel               float64
	maxNumberOfParticles     float64
	jitter                   float64
	isOneShot                bool
	allowOverlapping         bool
	seed                     int64
	isEnabled                bool
	pointsGen                *TrianglePointGenerator
	numberOfEmittedParticles float64
}

func NewVolumeParticleEmitter2(
	implicitSurface *ImplicitSurfaceSet2,
	maxRegion *BoundingBox2D,
	spacing float64,
	initialVel *Vector3D.Vector3D,
) *VolumeParticleEmitter2 {
	return &VolumeParticleEmitter2{
		implicitSurface:          implicitSurface,
		bounds:                   maxRegion,
		spacing:                  spacing,
		initialVel:               initialVel,
		linearVel:                Vector3D.NewVector(0, 0, 0),
		angularVel:               0,
		maxNumberOfParticles:     constants.KMaxSize,
		numberOfEmittedParticles: 0,
		jitter:                   0,
		isOneShot:                true,
		allowOverlapping:         false,
		seed:                     0,
		isEnabled:                true,
		pointsGen:                NewTrianglePointGenerator(),
	}
}

func (e *VolumeParticleEmitter2) setTarget(particles *SphSystemData2) {

	e.particles = particles
	e.onSetTarget(particles)
}

func (e *VolumeParticleEmitter2) onSetTarget(particles *SphSystemData2) {

	// Do nothing.
}

func (e *VolumeParticleEmitter2) onUpdate() {

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

func (e *VolumeParticleEmitter2) emit(particles *SphSystemData2, newPositions, newVelocities *[]*Vector3D.Vector3D) {

	e.implicitSurface.updateQueryEngine()

	region := NewBoundingBox2D(e.bounds.lowerCorner, e.bounds.upperCorner)

	if e.implicitSurface.isBounded() {
		//todo: surfaceBBox:=
	}

	// Reserving more space for jittering
	j := e.jitter
	maxJitterDist := 0.5 * j * e.spacing
	numNewParticles := 0.0

	callback := func(points *([]*Vector3D.Vector3D), point *Vector3D.Vector3D) bool {
		newAngleInRadian := (rand.Float64() - 0.5) * math.Pi * 2
		rotationMatrix := makeRotationMatrix(newAngleInRadian)
		randomDir := rotationMatrix.MultiplyMatrixByTuple(Vector3D.NewVector(0, 0, 0))
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

	if e.allowOverlapping || e.isOneShot {

		e.pointsGen.forEachPoint(region, e.spacing, nil, callback)
	} else {
		// Use serial hash grid searcher for continuous update.
		// todo.
	}
	// not needed?
	//*newVelocities = make([]*Vector3D.Vector3D, len(*newPositions), len(*newPositions))

	// original code from: \FluidEngine\fluid-engine-dev\src\jet\volume_particle_emitter2.cpp
	//newVelocities->parallelForEachIndex([&](size_t i) {
	//	(*newVelocities)[i] = velocityAt((*newPositions)[i]);
	//});

	e.parallelForEachIndex(newVelocities, newPositions)
}

func (e *VolumeParticleEmitter2) velocityAt(point *Vector3D.Vector3D) *Vector3D.Vector3D {

	r := point.Substract(e.implicitSurface.surfaces[0].getTransform().translation)
	a := Vector3D.NewVector(-r.Y, r.X, 0).Multiply(e.angularVel)
	return a.Add(e.linearVel).Add(e.initialVel)
}

func (e *VolumeParticleEmitter2) parallelForEachIndex(newVelocities, newPositions *[]*Vector3D.Vector3D) {
	for i := 0.0; i < e.numberOfEmittedParticles; i++ {

		e.callback(i, newVelocities, newPositions)
	}
}

func (e *VolumeParticleEmitter2) callback(i float64, newVelocities, newPositions *[]*Vector3D.Vector3D) {

	*newVelocities = append(*newVelocities, e.velocityAt((*newPositions)[int64(i)]))
}
