package dispersion

import (
	"fmt"

	"fiber-cd/internal/model"
)

func ProfileDeltaSlope(core, clad Sellmeier, lambdaUm, hUm float64) (float64, error) {
	if hUm <= 0 {
		return 0, fmt.Errorf("dispersion: finite-difference step must be > 0")
	}
	d1, err := DeltaOf(core, clad, lambdaUm+hUm)
	if err != nil {
		return 0, err
	}
	d0, err := DeltaOf(core, clad, lambdaUm-hUm)
	if err != nil {
		return 0, err
	}
	return (d1 - d0) / (2 * hUm), nil
}

func ProfileDispersion(core, clad Sellmeier, lambdaUm float64) (float64, error) {
	n1 := core.IndexUm(lambdaUm)
	delta, err := DeltaOf(core, clad, lambdaUm)
	if err != nil {
		return 0, err
	}
	slope, err := ProfileDeltaSlope(core, clad, lambdaUm, 0.002)
	if err != nil {
		return 0, err
	}
	return -(n1 * delta / SpeedOfLightMPS) * (lambdaUm * slope / delta) * 1e12, nil
}

func ComposeWithProfile(cfg model.Config, moleFrac float64) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	core, err := MixSellmeier(Silica(), Germania(), moleFrac)
	if err != nil {
		return Result{}, err
	}
	clad := Silica()
	lambdaUm := cfg.WavelengthUm()
	n1 := core.IndexUm(lambdaUm)
	n2 := clad.IndexUm(lambdaUm)
	if n1 <= n2 {
		return Result{}, fmt.Errorf("dispersion: doped core not guiding at %v nm", cfg.WavelengthNm)
	}
	na := NumericalAperture(n1, n2)
	v := vAt(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	dMat := MaterialDispersionOf(core, lambdaUm)
	dWg := WaveguideDispersion(n1, (n1-n2)/n1, cfg.WavelengthM(), v)
	dPr, err := ProfileDispersion(core, clad, lambdaUm)
	if err != nil {
		return Result{}, err
	}
	return Result{
		LambdaNm: cfg.WavelengthNm,
		V:        v,
		DMat:     dMat,
		DWg:      dWg,
		DTotal:   dMat + dWg + dPr,
	}, nil
}

func MaterialZeroShiftNm(moleFrac float64) (float64, error) {
	pure, err := MaterialZeroOf(Silica(), 1.20, 1.40)
	if err != nil {
		return 0, err
	}
	mix, err := MixSellmeier(Silica(), Germania(), moleFrac)
	if err != nil {
		return 0, err
	}
	doped, err := MaterialZeroOf(mix, 1.20, 1.50)
	if err != nil {
		return 0, err
	}
	return (doped - pure) * 1e3, nil
}

func WaveguideKeepsNegativeSign(cfg model.Config) (bool, error) {
	res, err := TotalDispersionAt(cfg)
	if err != nil {
		return false, err
	}
	return res.DWg < 0, nil
}
