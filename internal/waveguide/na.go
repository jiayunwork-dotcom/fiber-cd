package waveguide

import "math"

func NumericalAperture(n1, n2 float64) float64 {
	return math.Sqrt(n1*n1 - n2*n2)
}
