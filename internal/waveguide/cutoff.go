package waveguide

import "math"

const CutoffV = 2.405

func CutoffWavelengthM(coreRadiusM, na float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / CutoffV
}

func CutoffWavelengthNm(coreRadiusM, na float64) float64 {
	return CutoffWavelengthM(coreRadiusM, na) * 1e9
}
