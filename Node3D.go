package main

import "jimmykiang/fluidengine/constants"

type Node3D struct {
	flags string
	bound *BoundingBox3D
	child float64
	item  float64
}

func NewNode3D() *Node3D {
	return &Node3D{
		flags: "0",
		bound: NewBoundingBox3DReset(),
		child: constants.KMaxSize,
		item:  constants.KMaxSize,
	}
}

func (n *Node3D) initLeaf(it float64, b *BoundingBox3D) {
	n.flags = "3"
	n.item = it
	n.child = it
	n.bound = b
}
