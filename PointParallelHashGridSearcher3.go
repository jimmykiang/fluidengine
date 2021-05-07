package main

import "math"

// PointParallelHashGridSearcher3 is a parallel version of hash grid-based 3-D point searcher.
// This struct implements parallel version of 3-D point searcher by using hash
// grid for its internal acceleration data structure. Each point is recorded to
// its corresponding bucket where the hashing function is 3-D grid mapping.
type PointParallelHashGridSearcher3 struct {
	gridSpacing     float64
	resolution      *Vector3D
	points          []*Vector3D
	startIndexTable []int64
	endIndexTable   []int64
	sortedIndices   []int64
	keys            []int64
}

func NewPointParallelHashGridSearcher3(
	resolutionX float64,
	resolutionY float64,
	resolutionZ float64,
	gridSpacing float64,
) *PointParallelHashGridSearcher3 {
	return &PointParallelHashGridSearcher3{
		gridSpacing: gridSpacing,
		resolution: NewVector(
			math.Max(resolutionX, kOneSSize),
			math.Max(resolutionY, kOneSSize),
			math.Max(resolutionZ, kOneSSize),
		),
		points:          make([]*Vector3D, 5),
		startIndexTable: make([]int64, 5),
		endIndexTable:   make([]int64, 5),
		sortedIndices:   make([]int64, 5),
		keys:            make([]int64, 5),
	}
}
