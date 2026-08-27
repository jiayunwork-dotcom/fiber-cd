package waveguide

import "fiber-cd/internal/model"

func Analyze(cfg model.Config) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	na := NumericalAperture(cfg.N1, cfg.N2)
	delta := RelativeIndexDelta(cfg.N1, cfg.N2)
	v := VNumber(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	cutoffNm := CutoffWavelengthNm(cfg.CoreRadiusM(), na)
	return Result{
		Config:             cfg,
		NA:                 na,
		Delta:              delta,
		V:                  v,
		CutoffWavelengthNm: cutoffNm,
		Mode:               Classify(v),
	}, nil
}
