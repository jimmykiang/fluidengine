package physicsHelper

import "math"

func ComputePressureFromEos(
	density float64,
	targetDensity float64,
	eosScale float64,
	eosExponent float64,
	negativePressureScale float64,
) float64 {
	p := eosScale / eosExponent * (math.Pow(density/targetDensity, eosExponent) - 1)

	// Negative pressure scaling.
	if p < 0 {
		p *= negativePressureScale
	}
	return p
}
