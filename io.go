package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sbinet/npyio"
)

func saveNpy(path, conf, fileName string, results []float64, frame *Frame) {

	f, err := os.Create(filepath.Join(path, conf, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write to .npy with the history of past values.
	m := results[:frame.index]
	if frame.index == 0 {

		// must always pass a slice with the single element of interest when index == 0.
		m = []float64{results[frame.index]}
	}

	err = npyio.Write(f, m)
	if err != nil {
		log.Fatalf("error writing to file: %v\n", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v\n", err)
	}
}
