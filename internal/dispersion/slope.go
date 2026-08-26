package dispersion

import "fiber-cd/internal/model"

type SlopeResult struct {
	LambdaNm   float64
	STotal     float64
	SMaterial  float64
	SWaveguide float64
}

func dTotalAtNm(cfg model.Config, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, lambdaM)
	return dispersionAt(cfg, v, lambdaM).DTotal
}

func TotalDispersionSlope(cfg model.Config, lambdaNm float64) (float64, error) {
	if err := model.Validate(cfg); err != nil {
		return 0, err
	}
	const h = 0.5
	return (dTotalAtNm(cfg, lambdaNm+h) - dTotalAtNm(cfg, lambdaNm-h)) / (2 * h), nil
}

func SlopeAtOperating(cfg model.Config) (SlopeResult, error) {
	if err := model.Validate(cfg); err != nil {
		return SlopeResult{}, err
	}
	sTotal, err := TotalDispersionSlope(cfg, cfg.WavelengthNm)
	if err != nil {
		return SlopeResult{}, err
	}
	sMat := MaterialDispersionSlopeUm(cfg.WavelengthUm())
	return SlopeResult{
		LambdaNm:   cfg.WavelengthNm,
		STotal:     sTotal,
		SMaterial:  sMat,
		SWaveguide: sTotal - sMat,
	}, nil
}
