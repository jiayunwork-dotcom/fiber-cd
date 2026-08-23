package dispersion

import "fiber-cd/internal/model"

// SlopeResult 是一次色散斜率核算的产物，单位 ps/(nm²·km)。
type SlopeResult struct {
	// LambdaNm 核算波长（nm）。
	LambdaNm float64
	// STotal 总色散斜率 dD_tot/dλ。
	STotal float64
	// SMaterial 材料项斜率 dD_mat/dλ（闭合形式）。
	SMaterial float64
	// SWaveguide 波导项斜率 = STotal − SMaterial。
	SWaveguide float64
}

// dTotalAtNm 返回 D_tot 在任意波长（nm）下的取值。
func dTotalAtNm(cfg model.Config, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, lambdaM)
	return dispersionAt(cfg, v, lambdaM).DTotal
}

// TotalDispersionSlope 用中心差分求 D_tot 在 lambdaNm 处的导数：
//
//	S(λ) ≈ [D_tot(λ+h) − D_tot(λ−h)] / 2h，h = 0.5 nm
//
// 与闭合形式的材料斜率交叉验证时，二者应彼此吻合。
func TotalDispersionSlope(cfg model.Config, lambdaNm float64) (float64, error) {
	if err := model.Validate(cfg); err != nil {
		return 0, err
	}
	const h = 0.5 // nm
	return (dTotalAtNm(cfg, lambdaNm+h) - dTotalAtNm(cfg, lambdaNm-h)) / (2 * h), nil
}

// SlopeAtOperating 在配置的工作波长下合成总斜率、材料斜率与波导斜率。
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
