package main

import (
	"math"
	"math/rand"
)

// PointParticleEmitter3 represents a 3-D particle emitter.
//This class emits particles from a single point in given direction, speed, and spreading angle.
type PointParticleEmitter3 struct {
	isEnabled                        bool
	particles                        *ParticleSystemData3
	onUpdateCallback                 OnBeginUpdateCallbackEmitter
	rng                              int
	firstFrameTimeInSeconds          float64
	numberOfEmittedParticles         int
	maxNumberOfNewParticlesPerSecond int
	maxNumberOfParticles             uint64
	origin                           *Vector3D
	direction                        *Vector3D
	speed                            float64
	spreadAngleInRadians             float64
	seed                             uint32
}

// OnBeginUpdateCallbackEmitter is a brief Callback function signature type for update calls.
// This type of callback function will take the emitter pointer, current
// time, and time interval in seconds.
type OnBeginUpdateCallbackEmitter func(
	rigidBodyCollider *PointParticleEmitter3,
	currentTime float64,
	timeInterval float64,
)

func NewPointParticleEmitter3() *PointParticleEmitter3 {
	return &PointParticleEmitter3{
		isEnabled:                        true,
		particles:                        nil,
		onUpdateCallback:                 nil,
		rng:                              0,
		firstFrameTimeInSeconds:          0,
		numberOfEmittedParticles:         0,
		maxNumberOfNewParticlesPerSecond: 0,
		maxNumberOfParticles:             18446744073709551615,
		origin:                           NewVector(0, 0, 0),
		direction:                        NewVector(0, 1, 0),
		speed:                            1,
		spreadAngleInRadians:             0,
		seed:                             0,
	}
}

func (e *PointParticleEmitter3) withOrigin(v *Vector3D) {

	e.origin.Set(v)
}

func (e *PointParticleEmitter3) withDirection(v *Vector3D) {

	e.direction.Set(v)
}

func (e *PointParticleEmitter3) withSpeed(s float64) {

	e.speed = s
}

func (e *PointParticleEmitter3) withSpreadAngleInDegrees(d float64) {

	e.spreadAngleInRadians = degreesToRadians(d)
}

func (e *PointParticleEmitter3) withMaxNumberOfNewParticlesPerSecond(m int) {

	e.maxNumberOfNewParticlesPerSecond = m
}

func (e *PointParticleEmitter3) setTarget(particles *ParticleSystemData3) {

	e.particles = particles
	e.onSetTarget(particles)
}

func (e *PointParticleEmitter3) onSetTarget(particles *ParticleSystemData3) {

	// Do nothing.
}

func (e *PointParticleEmitter3) update(currentTimeInSeconds float64, timeIntervalInSeconds float64) {



	if e.particles == nil{
		return
	}

	particles := e.particles

	if e.numberOfEmittedParticles == 0{

		e.firstFrameTimeInSeconds = currentTimeInSeconds
	}

	elapsedTimeInSeconds := currentTimeInSeconds - e.firstFrameTimeInSeconds

	newMaxTotalNumberOfEmittedParticles := math.Ceil((elapsedTimeInSeconds+timeIntervalInSeconds) *
		float64(e.maxNumberOfNewParticlesPerSecond))

	newMaxTotalNumberOfEmittedParticles = math.Min(
		newMaxTotalNumberOfEmittedParticles,
		float64(e.maxNumberOfParticles),
	)

	maxNumberOfNewParticles := newMaxTotalNumberOfEmittedParticles - float64(e.numberOfEmittedParticles)

	if maxNumberOfNewParticles > 0 {

		candidatePositions := make([]*Vector3D, 0)
		candidateVelocities := make([]*Vector3D, 0)
		newPositions := make([]*Vector3D, 0)
		newVelocities := make([]*Vector3D, 0)

		e.emit(&candidatePositions, &candidateVelocities, maxNumberOfNewParticles)

		newPositions = append(newPositions, candidatePositions...)
		newVelocities = append(newVelocities, candidateVelocities...)

		particles.addParticles(newPositions, newVelocities, nil)

		e.numberOfEmittedParticles += len(newPositions)
	}
}

func (e *PointParticleEmitter3) emit(
	newPositions *[]*Vector3D,
	newVelocities *[]*Vector3D,
	maxNewNumberOfParticles float64,
) {

	for i := 0; i < int(maxNewNumberOfParticles); i++ {

		newDirection := e.uniformSampleCone(rand.Float64(), rand.Float64(), e.direction, e.spreadAngleInRadians)
		*newPositions = append(*newPositions, e.origin)
		*newVelocities = append(*newVelocities, newDirection.Multiply(e.speed))
	}
}

// uniformSampleCone returns randomly sampled direction within a cone.
// u1    First random sample.
// u2    Second random sample.
// axis  The axis of the cone.
// angle The angle of the cone.
// return     Sampled direction vector.
func (e *PointParticleEmitter3) uniformSampleCone(u1 float64, u2 float64, axis *Vector3D, angle float64) *Vector3D {

	cosAngle_2 := math.Cos(angle / 2)
	y := 1 - (1-cosAngle_2)*u1
	r := math.Sqrt(math.Max(0, 1-y*y))
	phi := math.Pi * 2 * u2
	x := r * math.Cos(phi)
	z := r * math.Sin(phi)
	ts := axis.tangential()

	a := ts[0].Multiply(x)
	b := axis.Multiply(y)
	c := ts[1].Multiply(z)

	return a.Add(b).Add(c)
}
