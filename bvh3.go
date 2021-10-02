package main

type Bvh3 struct {
	IntersectionQueryEngine3    *IntersectionQueryEngine3
	NearestNeighborQueryEngine3 *NearestNeighborQueryEngine3
	items                       []ImplicitSurface3
	bound                       *BoundingBox3D
	itemBounds                  []*BoundingBox3D
	nodes                       []*Node3D
}

func NewBvh3() *Bvh3 {
	return &Bvh3{nodes: make([]*Node3D, 0, 0)}
}

// build the bounding volume hierarchy.
func (b *Bvh3) build(items []ImplicitSurface3, itemsBounds []*BoundingBox3D) {

	b.items = items
	b.itemBounds = itemsBounds

	if len(items) == 0 {
		return
	}
	b.nodes = nil
	b.bound = NewBoundingBox3DReset()

	itemsize := float64(len(b.items))
	for i := float64(0); i < itemsize; i++ {
		b.bound.merge(b.itemBounds[int(i)])
	}

	itemIndices := make([]float64, 0, int64(itemsize))

	for i := float64(0); i < itemsize; i++ {

		itemIndices = append(itemIndices, i)
	}

	b.buildInternal(0, itemIndices, itemsize, 0)
}

func (b *Bvh3) buildInternal(nodeIndex int, itemIndices []float64, nItems float64, currentDepth int) int {

	// add a node.
	b.nodes = append(b.nodes, NewNode3D())

	// initialize leaf node if termination criteria met.
	if nItems == 1 {
		b.nodes[nodeIndex].initLeaf(itemIndices[0], b.itemBounds[int64(itemIndices[0])])
		return currentDepth + 1
	}

	return 0
}
