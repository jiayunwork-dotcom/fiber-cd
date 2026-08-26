package waveguide

import (
	"fmt"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
)

type ProbeResult struct {
	Config       model.Config
	WavelengthNm float64
	V            float64
	Mode         ModeStatus
	CutoffNm     float64
	Disp         dispersion.Result
}

func Probe(cfg model.Config, wavelengthNm float64) (ProbeResult, error) {
	if err := model.Validate(cfg); err != nil {
		return ProbeResult{}, err
	}
	if wavelengthNm <= 0 {
		return ProbeResult{}, fmt.Errorf("probe wavelength must be positive, got %g nm", wavelengthNm)
	}
	na := NumericalAperture(cfg.N1, cfg.N2)
	lambdaM := wavelengthNm * 1e-9
	v := VNumber(cfg.CoreRadiusM(), na, lambdaM)
	d, err := dispersion.Compose(cfg, v, lambdaM)
	if err != nil {
		return ProbeResult{}, err
	}
	return ProbeResult{
		Config:       cfg,
		WavelengthNm: wavelengthNm,
		V:            v,
		Mode:         Classify(v),
		CutoffNm:     CutoffWavelengthNm(cfg.CoreRadiusM(), na),
		Disp:         d,
	}, nil
}
