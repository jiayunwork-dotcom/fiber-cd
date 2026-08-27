package dispersion

import (
	"fmt"

	"fiber-cd/internal/model"
)

type SpectrumPoint struct {
	LambdaNm   float64
	LambdaUm   float64
	Index      float64
	GroupIndex float64
	DMat       float64
	DWg        float64
	DTotal     float64
	STotal     float64
}

func SpectrumAt(cfg model.Config, lambdaNm float64) (SpectrumPoint, error) {
	if err := model.Validate(cfg); err != nil {
		return SpectrumPoint{}, err
	}
	if lambdaNm <= 0 {
		return SpectrumPoint{}, fmt.Errorf("wavelength must be positive, got %g nm", lambdaNm)
	}
	lambdaM := lambdaNm * 1e-9
	lambdaUm := lambdaNm / 1000

	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, lambdaM)
	d := dispersionAt(cfg, v, lambdaM)
	s, err := TotalDispersionSlope(cfg, lambdaNm)
	if err != nil {
		return SpectrumPoint{}, err
	}

	silica := Silica()
	return SpectrumPoint{
		LambdaNm:   lambdaNm,
		LambdaUm:   lambdaUm,
		Index:      silica.IndexUm(lambdaUm),
		GroupIndex: GroupIndex(lambdaUm),
		DMat:       d.DMat,
		DWg:        d.DWg,
		DTotal:     d.DTotal,
		STotal:     s,
	}, nil
}
