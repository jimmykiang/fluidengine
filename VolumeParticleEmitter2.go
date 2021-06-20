package main

// VolumeParticleEmitter2 is a 2-D volumetric particle emitter.
// emits particles from volumetric geometry.
// Constructs an emitter that spawns particles from given implicit surface
// which defines the volumetric geometry. Provided bounding box limits
// the particle generation region.
type VolumeParticleEmitter2 struct {
	implicitSurface      *ImplicitSurfaceSet2
	particles            *SphSystemData2
	bounds               *BoundingBox2D
	spacing              float64
	initialVel           *Vector3D
	linearVel            *Vector3D
	angularVel           float64
	maxNumberOfParticles float64
	jitter               float64
	isOneShot            bool
	allowOverlapping     bool
	seed                 int64
	isEnabled            bool
	pointsGen            *TrianglePointGenerator
}

func NewVolumeParticleEmitter2(
	implicitSurface *ImplicitSurfaceSet2,
	maxRegion *BoundingBox2D,
	spacing float64,
	initialVel *Vector3D,
) *VolumeParticleEmitter2 {
	return &VolumeParticleEmitter2{
		implicitSurface:      implicitSurface,
		bounds:               maxRegion,
		spacing:              spacing,
		initialVel:           initialVel,
		linearVel:            NewVector(0, 0, 0),
		angularVel:           0,
		maxNumberOfParticles: kMaxSize,
		jitter:               0,
		isOneShot:            true,
		allowOverlapping:     false,
		seed:                 0,
		isEnabled:            true,
		pointsGen:            NewTrianglePointGenerator(),
	}
}

func (e *VolumeParticleEmitter2) setTarget(particles *SphSystemData2) {

	e.particles = particles
	e.onSetTarget(particles)
}

func (e *VolumeParticleEmitter2) onSetTarget(particles *SphSystemData2) {

	// Do nothing.
}
