package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/constants"
	"math"
)

// PointParallelHashGridSearcher3 is a parallel version of hash grid-based 3-D point searcher.
// This struct implements parallel version of 3-D point searcher by using hash
// grid for its internal acceleration data structure. Each point is recorded to
// its corresponding bucket where the hashing function is 3-D grid mapping.
type PointParallelHashGridSearcher3 struct {
	gridSpacing     float64
	resolution      *Vector3D.Vector3D
	points          []*Vector3D.Vector3D
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
		resolution: Vector3D.NewVector(
			math.Max(resolutionX, constants.KOneSSize),
			math.Max(resolutionY, constants.KOneSSize),
			math.Max(resolutionZ, constants.KOneSSize),
		),
		points:          make([]*Vector3D.Vector3D, 0, 0),
		startIndexTable: make([]int64, 0, 0),
		endIndexTable:   make([]int64, 0, 0),
		sortedIndices:   make([]int64, 0, 0),
		keys:            make([]int64, 0, 0),
	}
}

func (s *PointParallelHashGridSearcher3) build(points []*Vector3D.Vector3D) {

	// Allocate memory chuncks.
	numberOfPoints := len(points)
	tempKeys := make([]int64, 0, 0)

	for i := 0; i < numberOfPoints; i++ {
		tempKeys = append(tempKeys, 0)
	}

	s.startIndexTable = make([]int64, 0, 0)
	for i := 0; i < int(s.resolution.X*s.resolution.Y); i++ {
		s.startIndexTable = append(s.startIndexTable, math.MaxInt64)
	}

	s.endIndexTable = make([]int64, 0, 0)
	for i := 0; i < int(s.resolution.X*s.resolution.Y); i++ {
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

	s.points = make([]*Vector3D.Vector3D, 0, 0)
	for i := 0; i < numberOfPoints; i++ {
		s.points = append(s.points, Vector3D.NewVector(0, 0, 0))
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
	tempKeysResult := make([]int64, 0, 0)

	uniqueTempKeys := make([]int64, 0, 0)
	keys := make(map[int64]bool)

	for _, entry := range tempKeys {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueTempKeys = append(uniqueTempKeys, entry)
		}
	}

	for _, tV := range uniqueTempKeys {
		for k, v := range tempKeys {
			if v == tV {
				tempKeysResult = append(tempKeysResult, int64(k))
			}
		}
	}
	copy(s.sortedIndices, tempKeysResult)

	//sort.Slice(s.sortedIndices, func(i, j int) bool { return tempKeys[i] < tempKeys[i] })

	// Re-order point and key arrays.
	for i := 0; i < numberOfPoints; i++ {
		s.points[i] = points[s.sortedIndices[i]]
		s.keys[i] = tempKeys[s.sortedIndices[i]]
	}

	// Now _points and _keys are sorted by points' hash key values.
	// Let's fill in start/end index table with _keys.
	// Assume that _keys array looks like:
	// [5|8|8|10|10|10]
	// Then _startIndexTable and _endIndexTable should be like:
	// [.....|0|...|1|..|3|..]
	// [.....|1|...|3|..|6|..]
	//       ^5    ^8   ^10
	// So that _endIndexTable[i] - _startIndexTable[i] is the number points
	// in i-th table bucket.

	s.startIndexTable[s.keys[0]] = 0
	s.endIndexTable[s.keys[numberOfPoints-1]] = int64(numberOfPoints)

	for i := 1; i < numberOfPoints; i++ {
		if s.keys[i] > s.keys[i-1] {
			s.startIndexTable[s.keys[i]] = int64(i)
			s.endIndexTable[s.keys[i-1]] = int64(i)
		}
	}

	sumNumberOfPointsPerBucket := int64(0)
	maxNumberOfPointsPerBucket := int64(0)
	numberOfNonEmptyBucket := int64(0)

	for i := 0; i < len(s.startIndexTable); i++ {
		if s.startIndexTable[i] != math.MaxInt64 {
			numberOfPointsInBucket := s.endIndexTable[i] - s.startIndexTable[i]
			sumNumberOfPointsPerBucket += numberOfPointsInBucket
			maxNumberOfPointsPerBucket = int64(math.Max(
				float64(maxNumberOfPointsPerBucket),
				float64(numberOfPointsInBucket),
			))
			numberOfNonEmptyBucket++
		}
	}
}

func (s *PointParallelHashGridSearcher3) getHashKeyFromPosition(position *Vector3D.Vector3D) int64 {
	bucketIndex := s.getBucketIndex(position)

	return s.getHashKeyFromBucketIndex(bucketIndex)
}

func (s *PointParallelHashGridSearcher3) getBucketIndex(position *Vector3D.Vector3D) *Vector3D.Vector3D {
	bucketIndex := Vector3D.NewVector(0, 0, 0)
	bucketIndex.X = math.Floor(position.X / s.gridSpacing)
	bucketIndex.Y = math.Floor(position.Y / s.gridSpacing)
	return bucketIndex
}

func (s *PointParallelHashGridSearcher3) getHashKeyFromBucketIndex(bucketIndex *Vector3D.Vector3D) int64 {

	wrappedIndex := Vector3D.NewVector(bucketIndex.X, bucketIndex.Y, 0)
	wrappedIndex.X = float64(int64(bucketIndex.X) % int64(s.resolution.X))
	wrappedIndex.Y = float64(int64(bucketIndex.Y) % int64(s.resolution.Y))

	if wrappedIndex.X < 0 {
		wrappedIndex.X += s.resolution.X
	}
	if wrappedIndex.Y < 0 {
		wrappedIndex.Y += s.resolution.Y
	}

	return int64(wrappedIndex.Y*s.resolution.X + wrappedIndex.X)
}

func (s *PointParallelHashGridSearcher3) forEachNearbyPoint(
	origin *Vector3D.Vector3D,
	radius float64,
	iExternal int64,
	sum *float64,
	callback func(int64, int64, *Vector3D.Vector3D, *Vector3D.Vector3D, *float64),
) {
	nearbyKeys := make([]int64, 4, 4)
	s.getNearbyKeys(origin, nearbyKeys)

	queryRadiusSquared := radius * radius

	for i := 0; i < 4; i++ {
		nearbyKey := nearbyKeys[i]
		start := s.startIndexTable[nearbyKey]
		end := s.endIndexTable[nearbyKey]

		// Empty bucket -- continue to next bucket.
		if start == math.MaxInt64 {
			continue
		}

		for j := start; j < end; j++ {
			direction := s.points[j].Substract(origin)
			distanceSquared := direction.Squared()

			if distanceSquared <= queryRadiusSquared {
				callback(int64(iExternal), s.sortedIndices[j], s.points[j], origin, sum)
			}
		}
	}
}

func (s *PointParallelHashGridSearcher3) getNearbyKeys(
	position *Vector3D.Vector3D,
	nearbyKeys []int64,
) {
	originIndex := s.getBucketIndex(position)

	nearbyBucketIndices := make([]*Vector3D.Vector3D, 0, 0)

	for i := 0; i < 4; i++ {
		nearbyBucketIndices = append(nearbyBucketIndices, Vector3D.NewVector(originIndex.X, originIndex.Y, 0))
	}

	if ((originIndex.X + 0.5) * s.gridSpacing) <= position.X {
		nearbyBucketIndices[2].X += 1
		nearbyBucketIndices[3].X += 1
	} else {
		nearbyBucketIndices[2].X -= 1
		nearbyBucketIndices[3].X -= 1
	}

	if ((originIndex.Y + 0.5) * s.gridSpacing) <= position.Y {
		nearbyBucketIndices[1].Y += 1
		nearbyBucketIndices[3].Y += 1
	} else {
		nearbyBucketIndices[1].Y -= 1
		nearbyBucketIndices[3].Y -= 1
	}

	for i := 0; i < 4; i++ {
		nearbyKeys[i] = s.getHashKeyFromBucketIndex(nearbyBucketIndices[i])
	}
}
