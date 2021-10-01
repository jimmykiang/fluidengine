package main

import (
	"jimmykiang/fluidengine/Vector3D"
	"jimmykiang/fluidengine/mathHelper"
)

type Box3 struct {
	Surface3 *Surface3
	// Bounding box of this box.
	bound *BoundingBox3D
	// Local-to-world transform.
	transform       *Transform3
	isNormalFlipped bool
}

func NewBox3(boundingBox *BoundingBox3D) *Box3 {
	return &Box3{
		Surface3:        NewSurface3(),
		bound:           boundingBox,
		transform:       NewTransform3(),
		isNormalFlipped: true,
	}
}

// Returns the closest point from the given point otherPoint to the surface.
func (p *Box3) closestPoint(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.

	otherPointLocal := p.transform.toLocal(otherPoint)
	d := p.closestPointLocal(otherPointLocal)
	return p.transform.toWorld(d)
}

func (p *Box3) closestPointLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {
	if p.bound.contains(otherPoint) {
		planes := make([]*Plane2D, 0, 6)
		planes = append(planes, NewPlane2D(Vector3D.NewVector(1, 0, 0), p.bound.upperCorner))
		planes = append(planes, NewPlane2D(Vector3D.NewVector(0, 1, 0), p.bound.upperCorner))
		planes = append(planes, NewPlane2D(Vector3D.NewVector(0, 0, 1), p.bound.upperCorner))
		planes = append(planes, NewPlane2D(Vector3D.NewVector(-1, 0, 0), p.bound.lowerCorner))
		planes = append(planes, NewPlane2D(Vector3D.NewVector(0, -1, 0), p.bound.lowerCorner))
		planes = append(planes, NewPlane2D(Vector3D.NewVector(0, 0, -1), p.bound.lowerCorner))

		result := planes[0].closestPoint(otherPoint)
		distanceSquared := result.DistanceSquaredTo(otherPoint)

		for i := 1; i < 6; i++ {

			localResult := planes[i].closestPoint(otherPoint)
			localDistanceSquared := localResult.DistanceSquaredTo(otherPoint)

			if localDistanceSquared < distanceSquared {
				result = localResult
				distanceSquared = localDistanceSquared
			}
		}

		return result
	} else {
		return Vector3D.NewVector(mathHelper.Clamp(
			otherPoint.X,
			p.bound.lowerCorner.X,
			p.bound.upperCorner.X,
		), mathHelper.Clamp(
			otherPoint.Y,
			p.bound.lowerCorner.Y,
			p.bound.upperCorner.Y,
		), mathHelper.Clamp(
			otherPoint.Z,
			p.bound.lowerCorner.Z,
			p.bound.upperCorner.Z,
		))
	}
}

// Returns the normal to the closest point on the surface from the given otherPoint.
func (p *Box3) closestDistance(otherPoint *Vector3D.Vector3D) float64 {

	// Returns the closest distance from the given point otherPoint to the
	// point on the surface in local frame.
	otherPointLocal := p.transform.toLocal(otherPoint)
	d := p.closestPointLocal(otherPointLocal)

	//Returns the distance to the other vector.1

	return otherPointLocal.Substract(d).Length()
}

// closestNormal returns the normal to the closest point on the surface from the given
// point.
func (p *Box3) closestNormal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	result := p.transform.toWorldDirection(p.closestNormalLocal(otherPoint))
	if p.isNormalFlipped {

		//result.Multiply(-1)
		result = result.Multiply(-1)
	}

	return result
}

func (p *Box3) closestNormalLocal(otherPoint *Vector3D.Vector3D) *Vector3D.Vector3D {

	planes := make([]*Plane3D, 0, 6)
	planes = append(planes, NewPlane3D(Vector3D.NewVector(1, 0, 0), p.bound.upperCorner))
	planes = append(planes, NewPlane3D(Vector3D.NewVector(0, 1, 0), p.bound.upperCorner))
	planes = append(planes, NewPlane3D(Vector3D.NewVector(0, 0, 1), p.bound.upperCorner))
	planes = append(planes, NewPlane3D(Vector3D.NewVector(-1, 0, 0), p.bound.lowerCorner))
	planes = append(planes, NewPlane3D(Vector3D.NewVector(0, -1, 0), p.bound.lowerCorner))
	planes = append(planes, NewPlane3D(Vector3D.NewVector(0, 0, -1), p.bound.lowerCorner))

	if p.bound.contains(otherPoint) {
		closestNormal := planes[0].normal
		closestPoint := planes[0].closestPoint(otherPoint)
		minDistanceSquared := closestPoint.Substract(otherPoint).Squared()

		for i := 1; i < 6; i++ {

			localClosestPoint := planes[i].closestPoint(otherPoint)
			localDistanceSquared := localClosestPoint.Substract(otherPoint).Squared()

			if localDistanceSquared < minDistanceSquared {
				closestNormal = planes[i].normal
				minDistanceSquared = localDistanceSquared
			}
		}
		return closestNormal
	}
	return nil
}

func (p *Box3) getTransform() *Transform3 {
	return p.transform
}

// Returns true if otherPoint is inside the volume defined by the surface.
func (p *Box3) isInside(otherPoint *Vector3D.Vector3D) bool {

	return p.isNormalFlipped == !p.isInsideLocal(p.transform.toLocal(otherPoint))
}

// Returns true if otherPoint is inside by given depth the volume
// defined by the surface in local frame.
func (p *Box3) isInsideLocal(otherPointLocal *Vector3D.Vector3D) bool {

	cpLocal := p.closestPointLocal(otherPointLocal)
	normalLocal := p.closestNormalLocal(otherPointLocal)
	r := otherPointLocal.Substract(cpLocal)
	return r.DotProduct(normalLocal) < 0.0
}
