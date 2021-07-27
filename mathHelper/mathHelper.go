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

// Lerp      Computes linear interpolation.
// param[in]  f0    The first value.
// param[in]  f1    The second value.
// param[in]  t     Relative offset [0, 1] from the first value.
// tparam     S     Input value type.
// tparam     T     Offset type.
// return     The interpolated value.
func Lerp(value0, value1 *Vector3D.Vector3D, f float64) *Vector3D.Vector3D {

	a := 1 - f
	b := value0.Multiply(a)
	c := value1.Multiply(f)
	d := b.Add(c)

	return d
}
