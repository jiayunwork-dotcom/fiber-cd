package dispersion

import (
	"math"

	"fiber-cd/internal/model"
)

func NumericalAperture(n1, n2 float64) float64 {
	return math.Sqrt(n1*n1 - n2*n2)
}

func RelativeDelta(n1, n2 float64) float64 {
	return (n1 - n2) / n1
}

func vAt(coreRadiusM, na, lambdaM float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / lambdaM
}

func dispersionAt(cfg model.Config, v, lambdaM float64) Result {
	lambdaNm := lambdaM * 1e9
	lambdaUm := lambdaM * 1e6
	dMat := MaterialDispersionUm(lambdaUm)
	dWg := WaveguideDispersion(cfg.N1, RelativeDelta(cfg.N1, cfg.N2), lambdaM, v)
	return Result{
		LambdaNm: lambdaNm,
		V:        v,
		DMat:     dMat,
		DWg:      dWg,
		DTotal:   dMat + dWg,
	}
}

func Compose(cfg model.Config, v, lambdaM float64) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	return dispersionAt(cfg, v, lambdaM), nil
}

func TotalDispersionAt(cfg model.Config) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	return dispersionAt(cfg, v, cfg.WavelengthM()), nil
}
