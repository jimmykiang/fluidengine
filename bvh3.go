package main

type Bvh3 struct {
	IntersectionQueryEngine3    *IntersectionQueryEngine3
	NearestNeighborQueryEngine3 *NearestNeighborQueryEngine3
	items                       []ImplicitSurface3
	bound                       *BoundingBox3D
	itemBounds                  []*BoundingBox3D
	nodes                       []*Node
}

func NewBvh3() *Bvh3 {
	return &Bvh3{nodes: make([]*Node, 0, 0)}
}
