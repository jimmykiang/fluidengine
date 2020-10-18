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

	for frame := NewFrame(); frame.index < 240; frame.advance() {

		sineAnimation.onUpdate(frame)

		resultsX[frame.index] = frame.timeInSeconds()
		resultsY[frame.index] = sineAnimation.value

		path, err := os.Getwd()
		const conf = "animation/sineanimation"
		fileName := fmt.Sprintf("data.#line2,%04d,x.npy", frame.index)

		f, err := os.Create(filepath.Join(path, conf, fileName))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// write to .npy with the history of past values.
		m := resultsX[:frame.index]
		if frame.index == 0 {

			// must always pass a slice with the single element of interest when index == 0.
			m = []float64{resultsX[frame.index]}
		}

		err = npyio.Write(f, m)
		if err != nil {
			log.Fatalf("error writing to file: %v\n", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v\n", err)
		}

		path1, err1 := os.Getwd()
		const conf1 = "animation/sineanimation"

		fileName1 := fmt.Sprintf("data.#line2,%04d,y.npy", frame.index)

		f1, err1 := os.Create(filepath.Join(path1, conf1, fileName1))
		if err1 != nil {
			log.Fatal(err1)
		}
		defer f1.Close()

		m1 := resultsY[:frame.index]

		if frame.index == 0 {

			m1 = []float64{resultsY[frame.index]}
		}
		err1 = npyio.Write(f1, m1)
		if err != nil {
			log.Fatalf("error writing to file: %v\n", err1)
		}

		err1 = f1.Close()
		if err1 != nil {
			log.Fatalf("error closing file: %v\n", err1)
		}
	}

	path, err := os.Getwd()
	const conf = "animation/sineanimation"

	f, err := os.Create(filepath.Join(path, conf, "data.#line2,x.npy"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = npyio.Write(f, resultsX)
	if err != nil {
		log.Fatalf("error writing to file: %v\n", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v\n", err)
	}

	path1, err1 := os.Getwd()
	const conf1 = "animation/sineanimation"

	f1, err1 := os.Create(filepath.Join(path1, conf1, "data.#line2,y.npy"))
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f1.Close()

	err1 = npyio.Write(f1, resultsY)
	if err != nil {
		log.Fatalf("error writing to file: %v\n", err1)
	}

	err1 = f1.Close()
	if err1 != nil {
		log.Fatalf("error closing file: %v\n", err1)
	}
}
