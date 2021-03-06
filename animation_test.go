package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/sbinet/npyio"
)

// Test npyio library writing go slice data to npy format.
func TestWriteToNpyFile(t *testing.T) {

	path, err := os.Getwd()
	const conf = "animation/testwritetonpyfile"
	fileName := fmt.Sprintf("data.#line2,%04d,x.npy", 1)

	f, err := os.Create(filepath.Join(path, conf, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := []float64{0, 1, 2, 3, 4, 5}
	err = npyio.Write(f, m)
	if err != nil {
		log.Fatalf("error writing to file: %v\n", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v\n", err)
	}

	path1, err1 := os.Getwd()
	const conf1 = "animation/testwritetonpyfile"

	fileName1 := fmt.Sprintf("data.#line2,%04d,y.npy", 1)

	f1, err1 := os.Create(filepath.Join(path1, conf1, fileName1))
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f1.Close()

	m1 := []float64{0, 1, 2, 3, 4, 5}
	err1 = npyio.Write(f1, m1)
	if err != nil {
		log.Fatalf("error writing to file: %v\n", err1)
	}

	err1 = f1.Close()
	if err1 != nil {
		log.Fatalf("error closing file: %v\n", err1)
	}
}

func TestSineAnimation(t *testing.T) {

	resultsX := make([]float64, 240, 240)
	resultsY := make([]float64, 240, 240)

	sineAnimation := NewSineAnimation()
	frame := NewFrame()

	for ; frame.index < 240; frame.advance() {

		sineAnimation.onUpdate(frame)

		resultsX[frame.index] = frame.timeInSeconds()
		resultsY[frame.index] = sineAnimation.value

		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		const conf = "animation/sineanimation"
		fileNameX := fmt.Sprintf("data.#line2,%04d,x.npy", frame.index)
		fileNameY := fmt.Sprintf("data.#line2,%04d,y.npy", frame.index)

		saveNpy(path, conf, fileNameX, resultsX, frame)
		saveNpy(path, conf, fileNameY, resultsY, frame)
	}

	//path, err := os.Getwd()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//const conf = "animation/sineanimation"
	//
	//saveNpy(path, conf, "data.#line2,x.npy", resultsX, frame)
	//saveNpy(path, conf, "data.#line2,y.npy", resultsY, frame)
}

func TestSimpleMassSpringAnimation(t *testing.T) {

	var x []float64
	var y []float64

	anim := NewSimpleMassSpringAnimation()
	anim.makeChain(10)
	anim.wind = NewVector(-80.0, 190.0, 0.0)
	// Set constraint for pointIndex 13
	// and initial fixed position for "consistency" at x = -5 (because makeChain starts with index == 0 )
	// and avoid wipLash caused spring calculation in the simulation.
	anim.constraints = append(anim.constraints, &Constraint{6, NewVector(-5, 20, 0), NewVector(0, 0, 0)})
	anim.exportStates(&x, &y)

	frame := NewFrame()

	for ; frame.index < 1000; frame.advance() {

		if frame.index > 500 {
			anim.wind = NewVector(40.0, 30.0, 0.0)
		}
		anim.onUpdate(frame)
		anim.exportStates(&x, &y)

		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		const conf = "animation/simpleMassSpringAnimation"
		fileNameX := fmt.Sprintf("data.#line2,%04d,x.npy", frame.index)
		fileNameY := fmt.Sprintf("data.#line2,%04d,y.npy", frame.index)

		saveNpy(path, conf, fileNameX, x, frame)
		saveNpy(path, conf, fileNameY, y, frame)
	}
}

func TestParticleSystemSolver3HalfBounce(t *testing.T) {

	// Normal vector.
	normal1 := NewVector(0, 1, 0)
	// Point vector.
	point1 := NewVector(0, 0, 0)

	plane := NewPlane3D(normal1, point1)
	collider := NewRigidBodyCollider3(plane)
	solver := NewParticleSystemSolver3()
	emitter := NewPointParticleEmitter3()

	solver.SetCollider(collider)
	solver.SetDragCoefficient(0.0)
	solver.SetRestitutionCoefficient(0.5)
	solver.SetEmitter(emitter)

	particles := solver.ParticleSystemData()
	particles.addParticle(NewVector(0, 3, 0), NewVector(1, 0, 0), NewVector(0, 0, 0))

	x := make([]float64, 1000)
	y := make([]float64, 1000)

	frame := NewFrame()
	frame.timeIntervalInSeconds = 1.0 / 300.0

	for ; frame.index < 1000; frame.advance() {

		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)

		x[frame.index] = particles.positions()[0].x
		y[frame.index] = particles.positions()[0].y

		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		const conf = "animation/HalfBounce"
		fileNameX := fmt.Sprintf("data.#line2,%04d,x.npy", frame.index)
		fileNameY := fmt.Sprintf("data.#line2,%04d,y.npy", frame.index)

		saveNpy(path, conf, fileNameX, x, frame)
		saveNpy(path, conf, fileNameY, y, frame)
	}
}

func TestParticleSystemSolver3Update(t *testing.T) {

	// Normal vector.
	normal1 := NewVector(0, 1, 0)
	// Point vector.
	point1 := NewVector(0, 0, 0)
	plane := NewPlane3D(normal1, point1)
	collider := NewRigidBodyCollider3(plane)

	wind := NewConstantVectorField3()
	wind.withValue(NewVector(1, 0, 0))

	emitter := NewPointParticleEmitter3()
	emitter.withOrigin(NewVector(0, 3, 0))
	emitter.withDirection(NewVector(-1, 1, 0))
	emitter.withSpeed(5)
	emitter.withSpreadAngleInDegrees(45)
	emitter.withMaxNumberOfNewParticlesPerSecond(300)

	solver := NewParticleSystemSolver3()
	solver.SetCollider(collider)
	solver.SetEmitter(emitter)
	solver.setWind(wind)
	solver.SetDragCoefficient(0)
	solver.SetRestitutionCoefficient(0.5)

	frame := NewFrame()
	frame.timeIntervalInSeconds = 1.0 / 60.0

	ix := float64(-1)
	iy := float64(1)

	for ; frame.index < 500; frame.advance() {

		emitter.withDirection(NewVector(ix, iy, 0))

		// Some random wind
		//if frame.index < 200 {
		//	ix += 0.01
		//
		//}
		//if frame.index == 100 {
		//	iy =0
		//}
		//if frame.index > 200 {
		//	ix -= 0.04
		//	iy += 0.05
		//}

		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)

		solver.saveParticleDataXyUpdate(solver.particleSystemData, frame)
	}
}
