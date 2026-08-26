package waveguide

import "math"

const (
	MarcuseConstant = 0.65
	MarcuseMidTerm  = 1.619
	MarcuseTailTerm = 2.879
)

func SpotRadiusUm(coreRadiusUm, v float64) float64 {
	vv := math.Max(v, 1.0)
	ratio := MarcuseConstant + MarcuseMidTerm/math.Pow(vv, 1.5) + MarcuseTailTerm/math.Pow(vv, 6)
	return coreRadiusUm * ratio
}

func ModeFieldDiameterUm(coreRadiusUm, v float64) float64 {
	return 2 * SpotRadiusUm(coreRadiusUm, v)
}

func EffectiveAreaUm2(coreRadiusUm, v float64) float64 {
	w0 := SpotRadiusUm(coreRadiusUm, v)
	return math.Pi * w0 * w0
}
