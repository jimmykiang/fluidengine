package main

import (
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
	initialVel               *Vector3D
	linearVel                *Vector3D
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
	initialVel *Vector3D,
) *VolumeParticleEmitter2 {
	return &VolumeParticleEmitter2{
		implicitSurface:          implicitSurface,
		bounds:                   maxRegion,
		spacing:                  spacing,
		initialVel:               initialVel,
		linearVel:                NewVector(0, 0, 0),
		angularVel:               0,
		maxNumberOfParticles:     kMaxSize,
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

	newPositions := NewVector(0, 0, 0)
	newVelocities := NewVector(0, 0, 0)

	e.emit(particles, newPositions, newVelocities)
}

func (e *VolumeParticleEmitter2) emit(particles *SphSystemData2, newPositions *Vector3D, newVelocities *Vector3D) {

	e.implicitSurface.updateQueryEngine()

	region := NewBoundingBox2D(e.bounds.lowerCorner, e.bounds.upperCorner)

	if e.implicitSurface.isBounded() {
		//todo: surfaceBBox:=
	}

	// Reserving more space for jittering
	j := e.jitter
	maxJitterDist := 0.5 * j * e.spacing

	callback := func(points *([]*Vector3D), point *Vector3D) bool {
		newAngleInRadian := (rand.Float64() - 0.5) * math.Pi * 2
		rotationMatrix := makeRotationMatrix(newAngleInRadian)
		randomDir := rotationMatrix.MultiplyMatrixByTuple(NewVector(0, 0, 0))
		offset := randomDir.Multiply(maxJitterDist)
		candidate := point.Add(offset)

		if e.implicitSurface.signedDistance(candidate) <= 0.0 {
			if e.numberOfEmittedParticles < e.maxNumberOfParticles {
				// todo

			}
		}

		_, _ = randomDir, candidate
		return true
	}

	if e.allowOverlapping || e.isOneShot {

		e.pointsGen.forEachPoint(region, e.spacing, nil, callback)
	}
	_, _, _ = region, j, maxJitterDist
}
