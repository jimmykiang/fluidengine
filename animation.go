package main

import "math"

// Animation represents the base interface for the animation logic in its base level.
type Animation interface {

	// onUpdate should be overriden by downstream structs and implement its logic for updating the animation state.
	onUpdate(*Frame)
}

// SineAnimation contains the evaluated value for a typical sinusoid.
type SineAnimation struct {
	value float64
}

// NewSineAnimation creates and returns a new SineAnimation reference.
func NewSineAnimation() *SineAnimation {
	sineAnimation := &SineAnimation{
		value: 0,
	}

	return sineAnimation
}

// onUpdate for a standard sinusoidal function.
func (sineAnimation *SineAnimation) onUpdate(frame *Frame) {

	sineAnimation.value = math.Sin(10.0 * frame.timeInSeconds())
}

// Edge between 2 Vector3D points.
type Edge struct {
	first, second int
}

// NewEdge creates a new reference of Edge.
func NewEdge(first, second int) *Edge {
	return &Edge{first, second}
}

// Constraint by fixing the position of a point.
type Constraint struct {
	pointIndex    int
	fixedPosition *Vector3D
	fixedVelocity *Vector3D
}

// SimpleMassSpringAnimation contains the data for a mass-spring animation solver.
type SimpleMassSpringAnimation struct {
	positions, velocities, forces                           []*Vector3D
	edges                                                   []*Edge
	mass, stiffness, restLength, dampingCoefficient         float64
	dragCoefficient, floorPositionY, restitutionCoefficient float64
	gravity, wind                                           *Vector3D
	constraints                                             []*Constraint
}

// NewSimpleMassSpringAnimation creates and returns a new SimpleMassSpringAnimation reference.
func NewSimpleMassSpringAnimation() *SimpleMassSpringAnimation {
	simpleMassSpringAnimation := &SimpleMassSpringAnimation{
		mass:                   1.0,
		gravity:                NewVector(0.0, -9.8, 0.0),
		stiffness:              100.0,
		restLength:             1.0,
		dampingCoefficient:     1.0,
		dragCoefficient:        0.1,
		floorPositionY:         -7.0,
		restitutionCoefficient: 0.3,
	}

	return simpleMassSpringAnimation
}

// makeChain initializes the data by chaining the points horizontally.
func (anim *SimpleMassSpringAnimation) makeChain(numberOfPoints int) {

	if numberOfPoints == 0 {
		return
	}

	numberOfEdges := numberOfPoints - 1
	anim.positions = make([]*Vector3D, numberOfPoints)
	anim.velocities = make([]*Vector3D, numberOfPoints)
	anim.forces = make([]*Vector3D, numberOfPoints)

	for x := 0; x < numberOfPoints; x++ {
		anim.velocities[x] = NewVector(0, 0, 0)
		anim.forces[x] = NewVector(0, 0, 0)
	}

	anim.edges = make([]*Edge, numberOfEdges)

	for i := 0; i < numberOfPoints; i++ {

		anim.positions[i] = NewVector(-float64(i), 20, 0)
	}

	for i := 0; i < numberOfEdges; i++ {

		anim.edges[i] = NewEdge(i, i+1)
	}
}

// exportStates initializes the data by chaining the points horizontally.
func (anim *SimpleMassSpringAnimation) exportStates(x *[]float64, y *[]float64) {

	*x = make([]float64, len(anim.positions))
	*y = make([]float64, len(anim.positions))

	for i := 0; i < len(anim.positions); i++ {

		(*x)[i] = anim.positions[i].x
		(*y)[i] = anim.positions[i].y
	}
}

// onUpdate for SimpleMassSpringAnimation.
func (anim *SimpleMassSpringAnimation) onUpdate(frame *Frame) {

	numberOfPoints := len(anim.positions)
	numberOfEdges := len(anim.edges)

	// Compute forces.
	for i := 0; i < numberOfPoints; i++ {

		// Gravity force.
		anim.forces[i] = anim.gravity.Multiply(anim.mass)

		// Air drag force.
		relativeVel := anim.velocities[i]

		if anim.wind != nil {

			relativeVel = relativeVel.Substract(anim.wind)
		}
		anim.forces[i] = anim.forces[i].Substract(relativeVel.Multiply(anim.dragCoefficient))
	}

	for i := 0; i < numberOfEdges; i++ {

		pointIndex0 := anim.edges[i].first
		pointIndex1 := anim.edges[i].second

		// Compute spring force.
		pos0 := anim.positions[pointIndex0]
		pos1 := anim.positions[pointIndex1]
		r := pos0.Substract(pos1)

		distance := r.Length()

		if distance > 0.0 {

			force := r.Normalize().Multiply(-anim.stiffness * (distance - anim.restLength))
			anim.forces[pointIndex0] = anim.forces[pointIndex0].Add(force)
			anim.forces[pointIndex1] = anim.forces[pointIndex1].Substract(force)
		}

		// Add damping force.
		vel0 := anim.velocities[pointIndex0]
		vel1 := anim.velocities[pointIndex1]
		relativeVel0 := vel0.Substract(vel1)
		damping := relativeVel0.Multiply(-anim.dampingCoefficient)
		anim.forces[pointIndex0] = anim.forces[pointIndex0].Add(damping)
		anim.forces[pointIndex1] = anim.forces[pointIndex1].Substract(damping)
	}

	// Update states.
	for i := 0; i < numberOfPoints; i++ {

		// Compute new states.
		newAcceleration := anim.forces[i].Divide(anim.mass)
		newVelocity := anim.velocities[i].Add(newAcceleration.Multiply(frame.timeIntervalInSeconds))
		newPosition := anim.positions[i].Add(newVelocity.Multiply(frame.timeIntervalInSeconds))

		// Collision.
		if newPosition.y < anim.floorPositionY {

			newPosition.y = anim.floorPositionY
			if newVelocity.y < 0.0 {

				newVelocity.y *= -anim.restitutionCoefficient
				newPosition.y += frame.timeIntervalInSeconds * newVelocity.y
			}
		}

		// Update states.
		anim.velocities[i] = newVelocity
		anim.positions[i] = newPosition
	}

	// Apply constraints
	for i := 0; i < len(anim.constraints); i++ {

		pointIndex := anim.constraints[i].pointIndex
		//anim.positions[pointIndex] = anim.constraints[pointIndex].fixedPosition
		//anim.velocities[pointIndex] = anim.constraints[pointIndex].fixedVelocity

		// Fix position + velocity based one the constraint[0].
		anim.positions[pointIndex] = anim.constraints[0].fixedPosition
		anim.velocities[pointIndex] = anim.constraints[0].fixedVelocity
	}
}
