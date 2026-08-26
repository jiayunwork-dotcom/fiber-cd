package waveguide

import "math"

func VNumber(coreRadiusM, na, wavelengthM float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / wavelengthM
}
