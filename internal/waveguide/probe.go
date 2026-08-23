package waveguide

import (
	"fmt"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
)

// ProbeResult 是 probe 子命令的产物：同一根光纤在指定波长下的状态，
// 用于回答"换到别的波长还是不是单模"这类问题。
type ProbeResult struct {
	Config       model.Config
	WavelengthNm float64
	V            float64
	Mode         ModeStatus
	CutoffNm     float64
	Disp         dispersion.Result
}

// Probe 用配置的几何参数核算光纤在指定波长（nm）下的 V、模式与色散。
// 与 Analyze 的区别仅在于波长由参数给出，几何仍取自配置。
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
