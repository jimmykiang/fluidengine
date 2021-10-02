package main

import (
	"fmt"
	"jimmykiang/fluidengine/Vector3D"
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
	anim.wind = Vector3D.NewVector(-80.0, 190.0, 0.0)
	// Set constraint for pointIndex 13
	// and initial fixed position for "consistency" at x = -5 (because makeChain starts with index == 0 )
	// and avoid wipLash caused spring calculation in the simulation.
	anim.constraints = append(anim.constraints, &Constraint{6, Vector3D.NewVector(-5, 20, 0), Vector3D.NewVector(0, 0, 0)})
	anim.exportStates(&x, &y)

	frame := NewFrame()

	for ; frame.index < 1000; frame.advance() {

		if frame.index > 500 {
			anim.wind = Vector3D.NewVector(40.0, 30.0, 0.0)
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
	normal1 := Vector3D.NewVector(0, 1, 0)
	// Point vector.
	point1 := Vector3D.NewVector(0, 0, 0)

	plane := NewPlane3D(normal1, point1)
	collider := NewRigidBodyCollider3(plane)
	solver := NewParticleSystemSolver3()
	emitter := NewPointParticleEmitter3()

	solver.SetCollider(collider)
	solver.SetDragCoefficient(0.0)
	solver.SetRestitutionCoefficient(0.5)
	solver.SetEmitter(emitter)

	particles := solver.ParticleSystemData()
	particles.addParticle(Vector3D.NewVector(0, 3, 0), Vector3D.NewVector(1, 0, 0), Vector3D.NewVector(0, 0, 0))

	x := make([]float64, 1000)
	y := make([]float64, 1000)

	frame := NewFrame()
	frame.timeIntervalInSeconds = 1.0 / 300.0

	for ; frame.index < 1000; frame.advance() {

		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)

		x[frame.index] = particles.positions()[0].X
		y[frame.index] = particles.positions()[0].Y

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
	normal1 := Vector3D.NewVector(0, 1, 0)
	// Point vector.
	point1 := Vector3D.NewVector(0, 0, 0)
	plane := NewPlane3D(normal1, point1)
	collider := NewRigidBodyCollider3(plane)

	wind := NewConstantVectorField3()
	wind.withValue(Vector3D.NewVector(1, 0, 0))

	emitter := NewPointParticleEmitter3()
	emitter.withOrigin(Vector3D.NewVector(0, 3, 0))
	emitter.withDirection(Vector3D.NewVector(-1, 1, 0))
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

		emitter.withDirection(Vector3D.NewVector(ix, iy, 0))

		// Some random wind
		//if frame.index < 200 {
		//	ix += 0.01
		//
		//}
		//if frame.index == 100 {
		//	iy = 0
		//}
		//if frame.index > 200 {
		//	ix -= 0.04
		//	iy += 0.05
		//}

		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)

		solver.saveParticleDataXyUpdate(solver.particleSystemData, frame)

		// unComment to enable g3n openGl visualizer @ frame 100.
		//if frame.index == 100 {
		//	n := solver.particleSystemData.numberOfParticles
		//	visualizer.Visualize(solver.particleSystemData.positions(), n)
		//}
	}
}

func TestSphSolver2WaterDrop(t *testing.T) {

	targetSpacing := 0.02
	domain := NewBoundingBox2D(Vector3D.NewVector(0, 0, 0), Vector3D.NewVector(1, 2, 0))

	// Initialize solvers.
	solver := NewSphSolver2()
	solver.setPseudoViscosityCoefficient(0)

	particles := solver.particleSystemData
	particles.setTargetDensity(1000)
	particles.setTargetSpacing(targetSpacing)

	// Initialize source.
	surfaceSet := NewImplicitSurfaceSet2()
	v1 := Vector3D.NewVector(0, 1, 0)
	v2 := Vector3D.NewVector(0, 0.25*domain.height(), 0)
	p := NewPlane2D(v1, v2)
	surfaceSet.addExplicitSurface(p)

	s := NewSphere2(domain.midPoint(), domain.width()*0.15)
	surfaceSet.addExplicitSurface(s)

	sourceBound := NewBoundingBox2DFromStruct(domain)
	sourceBound.expand(-targetSpacing)

	emitter := NewVolumeParticleEmitter2(surfaceSet, sourceBound, targetSpacing, Vector3D.NewVector(0, 0, 0))

	solver.setEmitter(emitter)

	// Initialize boundary
	box := NewBox2(domain)
	box.Surface2.isNormalFlipped = true

	collider := NewRigidBodyCollider2(box)
	solver.setCollider(collider)

	frame := NewFrame()
	frame.timeIntervalInSeconds = 1.0 / 60.0

	for ; frame.index < 120; frame.advance() {
		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)
		solver.saveParticleDataXyUpdate(solver.particleSystemData.particleSystemData, frame)
	}
}

func TestSphSolver3WaterDrop(t *testing.T) {

	targetSpacing := 0.02
	domain := NewBoundingBox3D(Vector3D.NewVector(0, 0, 0), Vector3D.NewVector(1, 2, 1))

	// Initialize solvers.
	solver := NewSphSolver3()
	solver.setPseudoViscosityCoefficient(10.0)

	particles := solver.particleSystemData
	particles.setTargetDensity(1000)
	particles.setTargetSpacing(targetSpacing)

	// Initialize source.
	surfaceSet := NewImplicitSurfaceSet3()
	v1 := Vector3D.NewVector(0, 1, 0)
	v2 := Vector3D.NewVector(0, 0.25*domain.height(), 0)
	p := NewPlane3D(v1, v2)
	surfaceSet.addExplicitSurface(p)

	s := NewSphere3(domain.midPoint(), domain.width()*0.15)
	surfaceSet.addExplicitSurface(s)

	sourceBound := NewBoundingBox3DFromStruct(domain)
	sourceBound.expand(-targetSpacing)

	emitter := NewVolumeParticleEmitter3(surfaceSet, sourceBound, targetSpacing, Vector3D.NewVector(0, 0, 0))
	solver.setEmitter(emitter)

	// Initialize boundary
	box := NewBox3(domain)
	box.Surface3.isNormalFlipped = true

	collider := NewRigidBodyCollider3(box)
	solver.setCollider(collider)

	solver.setViscosityCoefficient(0.1)

	frame := NewFrame()
	frame.timeIntervalInSeconds = 1.0 / 60.0
	for ; frame.index < 120; frame.advance() {

		fmt.Println("Frame index:", frame.index)
		solver.onUpdate(frame)
	}
}
