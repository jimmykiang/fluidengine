package main

import "jimmykiang/fluidengine/constants"

type Node struct {
	flags string
	bound *BoundingBox2D
	child float64
	item  float64
}

func NewNode() *Node {
	return &Node{
		flags: "0",
		bound: NewBoundingBox2DReset(),
		child: constants.KMaxSize,
		item:  constants.KMaxSize,
	}
}

func (n *Node) initLeaf(it float64, b *BoundingBox2D) {
	n.flags = "2"
	n.item = it
	n.child = it
	n.bound = b
}
