package main

type Bvh2 struct {
	IntersectionQueryEngine2    *IntersectionQueryEngine2
	NearestNeighborQueryEngine2 *NearestNeighborQueryEngine2
	items                       []ImplicitSurface2
	bound                       *BoundingBox2D
	itemBounds                  []*BoundingBox2D
	nodes                       []*Node
}

func NewBvh2() *Bvh2 {
	return &Bvh2{nodes: make([]*Node, 0, 0)}
}

// build the bounding volume hierarchy.
func (b *Bvh2) build(items []ImplicitSurface2, itemsBounds []*BoundingBox2D) {

	b.items = items
	b.itemBounds = itemsBounds

	if len(items) == 0 {
		return
	}
	b.nodes = nil
	b.bound = NewBoundingBox2DReset()

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

func (b *Bvh2) buildInternal(nodeIndex int, itemIndices []float64, nItems float64, currentDepth int) int {

	// add a node.
	b.nodes = append(b.nodes, NewNode())

	// initialize leaf node if termination criteria met.
	if nItems == 1 {
		b.nodes[nodeIndex].initLeaf(itemIndices[0], b.itemBounds[int64(itemIndices[0])])
		return currentDepth + 1
	}

	return 0
}
