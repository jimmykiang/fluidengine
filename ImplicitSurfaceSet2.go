package main

// ImplicitSurfaceSet2 represents 2-D implicit surface set.
type ImplicitSurfaceSet2 struct {
	ImplicitSurface2  *ImplicitSurface2
	bvhInvalidated    bool
	surfaces          []*ImplicitSurface2
	unboundedSurfaces []*ImplicitSurface2
	bvh               *Bvh2
}

func NewImplicitSurfaceSet2() *ImplicitSurfaceSet2 {
	return &ImplicitSurfaceSet2{
		ImplicitSurface2: NewImplicitSurface2(),
		bvhInvalidated:   true,
		//surfaces,
		//unboundedSurfaces,
		//bvh,
	}
}
