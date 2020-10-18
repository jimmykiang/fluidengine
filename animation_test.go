package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/sbinet/npyio"
)

func TestWriteToNpyFile(t *testing.T) {

	path, err := os.Getwd()
	const conf = "animation/data.#line2,0000,x.npy"

	f, err := os.Create(filepath.Join(path, conf))
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
	const conf1 = "animation/data.#line2,0000,y.npy"

	f1, err1 := os.Create(filepath.Join(path1, conf1))
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
