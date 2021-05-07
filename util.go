package main

// min returns the smallest value from the slice.
func min(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < c {
			c = values[i]
		}
	}
	return c
}
