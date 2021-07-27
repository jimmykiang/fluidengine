package mathHelper

func Clamp(val, low, high float64) float64 {
	if val < low {
		return low
	} else if val > high {
		return high
	} else {
		return val
	}
}
