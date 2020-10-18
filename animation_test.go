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
	fileName := fmt.Sprintf("a data.#line2,%04d,x.npy", 1)

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

	fileName1 := fmt.Sprintf("a data.#line2,%04d,y.npy", 1)

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

}
