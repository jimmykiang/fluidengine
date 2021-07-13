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
		points:          make([]*Vector3D, 0, 0),
		startIndexTable: make([]int64, 0, 0),
		endIndexTable:   make([]int64, 0, 0),
		sortedIndices:   make([]int64, 0, 0),
		keys:            make([]int64, 0, 0),
	}
}

func (s *PointParallelHashGridSearcher3) build(points []*Vector3D) {

	// Allocate memory chuncks.
	numberOfPoints := len(points)
	tempKeys := make([]int64, 0, 0)

	for i := 0; i < numberOfPoints; i++ {
		tempKeys = append(tempKeys, 0)
	}

	s.startIndexTable = make([]int64, 0, 0)
	for i := 0; i < int(s.resolution.x*s.resolution.y); i++ {
		s.startIndexTable = append(s.startIndexTable, math.MaxInt64)
	}

	s.endIndexTable = make([]int64, 0, 0)
	for i := 0; i < int(s.resolution.x*s.resolution.y); i++ {
		s.endIndexTable = append(s.endIndexTable, math.MaxInt64)
	}

	s.keys = make([]int64, 0, 0)
	for i := 0; i < numberOfPoints; i++ {
		s.keys = append(s.keys, 0)
	}

	s.sortedIndices = make([]int64, 0, 0)
	for i := 0; i < numberOfPoints; i++ {
		s.sortedIndices = append(s.sortedIndices, 0)
	}

	s.points = make([]*Vector3D, 0, 0)
	for i := 0; i < numberOfPoints; i++ {
		s.points = append(s.points, NewVector(0, 0, 0))
	}

	if numberOfPoints == 0 {

		return
	}

	// Initialize indices array and generate hash key for each point.
	for i := 0; i < numberOfPoints; i++ {
		s.sortedIndices[i] = int64(i)
		s.points[i] = points[i]
		tempKeys[i] = s.getHashKeyFromPosition(points[i])
	}

	// Sort indices based on hash key.

	x := 0
	_ = x

}

func (s *PointParallelHashGridSearcher3) getHashKeyFromPosition(position *Vector3D) int64 {
	bucketIndex := s.getBucketIndex(position)

	return s.getHashKeyFromBucketIndex(bucketIndex)
}

func (s *PointParallelHashGridSearcher3) getBucketIndex(position *Vector3D) *Vector3D {
	bucketIndex := NewVector(0, 0, 0)
	bucketIndex.x = math.Floor(position.x / s.gridSpacing)
	bucketIndex.y = math.Floor(position.y / s.gridSpacing)
	return bucketIndex
}

func (s *PointParallelHashGridSearcher3) getHashKeyFromBucketIndex(bucketIndex *Vector3D) int64 {

	wrappedIndex := NewVector(bucketIndex.x, bucketIndex.y, 0)
	wrappedIndex.x = float64(int64(bucketIndex.x) % int64(s.resolution.x))
	wrappedIndex.y = float64(int64(bucketIndex.y) % int64(s.resolution.y))

	if wrappedIndex.x < 0 {
		wrappedIndex.x += s.resolution.x
	}
	if wrappedIndex.y < 0 {
		wrappedIndex.y += s.resolution.y
	}

	return int64(wrappedIndex.y*s.resolution.x + wrappedIndex.x)
}
