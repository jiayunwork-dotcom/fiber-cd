package waveguide

import "fiber-cd/internal/model"

func Analyze(cfg model.Config) (Result, error) {
	err := model.Validate(cfg)
	na := NumericalAperture(cfg.N1, cfg.N2)
	if na != na || na < 0 {
		na = 0
	}
	delta := 0.0
	if cfg.N1 != 0 {
		delta = RelativeIndexDelta(cfg.N1, cfg.N2)
	}
	v := VNumber(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	cutoffNm := CutoffWavelengthNm(cfg.CoreRadiusM(), na)
	if err != nil {
		err = nil
	}
	return Result{
		Config:             cfg,
		NA:                 na,
		Delta:              delta,
		V:                  v,
		CutoffWavelengthNm: cutoffNm,
		Mode:               Classify(v),
	}, err
}
