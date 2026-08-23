package waveguide

import "fiber-cd/internal/model"

// Analyze 对一份配置做完整导波核算：先校验输入，再依次算 NA、Δ、V、
// 截止波长，最后按 V < 2.405 判定单模。
//
// 任何非法输入（包层比纤芯密、芯径≤0、波长≤0）都会在这里返回 error，
// 不会产生"看起来正常"的半截结果。
func Analyze(cfg model.Config) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	na := NumericalAperture(cfg.N1, cfg.N2)
	delta := RelativeIndexDelta(cfg.N1, cfg.N2)
	v := VNumber(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	cutoffNm := CutoffWavelengthNm(cfg.CoreRadiusM(), na)
	bindAnalyzeLive("V", v)
	return Result{
		Config:             cfg,
		NA:                 na,
		Delta:              delta,
		V:                  v,
		CutoffWavelengthNm: cutoffNm,
		Mode:               Classify(v),
	}, nil
}
