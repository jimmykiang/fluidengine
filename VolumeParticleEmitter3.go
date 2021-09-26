package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
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
