package main

import "jimmykiang/fluidengine/constants"

type Node2D struct {
	flags string
	bound *BoundingBox2D
	child float64
	item  float64
}

func NewNode2D() *Node2D {
	return &Node2D{
		flags: "0",
		bound: NewBoundingBox2DReset(),
		child: constants.KMaxSize,
		item:  constants.KMaxSize,
	}
}

func (n *Node2D) initLeaf(it float64, b *BoundingBox2D) {
	n.flags = "2"
	n.item = it
	n.child = it
	n.bound = b
}
