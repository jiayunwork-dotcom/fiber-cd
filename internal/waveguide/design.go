package waveguide

import "math"

func MaxSingleModeDiameterUm(na, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	aMax := CutoffV * lambdaM / (2 * math.Pi * na)
	return aMax * 1e6 * 2
}

func MaxSingleModeNA(coreDiameterUm, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	dM := coreDiameterUm * 1e-6
	return CutoffV * lambdaM / (math.Pi * dM)
}

func SingleModeMargin(v float64) float64 {
	return CutoffV - v
}

func SingleModeMarginPercent(v float64) float64 {
	return SingleModeMargin(v) / CutoffV * 100
}
