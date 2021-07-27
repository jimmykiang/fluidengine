package mathHelper

import "jimmykiang/fluidengine/Vector3D"

func Clamp(val, low, high float64) float64 {
	if val < low {
		return low
	} else if val > high {
		return high
	} else {
		return val
	}
}

func lerp(value0, value1, high *Vector3D.Vector3D) float64 {

	return 0
}
