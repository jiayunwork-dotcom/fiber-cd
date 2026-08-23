package waveguide

import "math"

// 逆向设计量：给定期望单模的波长与 NA/芯径之一，反推另一个参数的
// 上限，帮助回答"什么配置能保持单模"。

// MaxSingleModeDiameterUm 在给定 NA 与波长（nm）下，保持单模
// （V ≤ 2.405）允许的最大芯径（μm）。由 V=2.405 反解：
//
//	d_max = 2·2.405·λ / (2π·NA)
func MaxSingleModeDiameterUm(na, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	aMax := CutoffV * lambdaM / (2 * math.Pi * na)
	return aMax * 1e6 * 2
}

// MaxSingleModeNA 在给定芯径（μm）与波长（nm）下，保持单模允许的
// 最大 NA：
//
//	NA_max = 2.405·λ / (π·d)
func MaxSingleModeNA(coreDiameterUm, lambdaNm float64) float64 {
	lambdaM := lambdaNm * 1e-9
	dM := coreDiameterUm * 1e-6
	return CutoffV * lambdaM / (math.Pi * dM)
}

// SingleModeMargin 返回单模裕量 2.405 − V：正值表示单模（越大越稳），
// 负值表示已进入多模区。
func SingleModeMargin(v float64) float64 {
	return CutoffV - v
}

// SingleModeMarginPercent 返回以 V 计的单模裕量百分比
// (2.405 − V)/2.405×100。
func SingleModeMarginPercent(v float64) float64 {
	return SingleModeMargin(v) / CutoffV * 100
}
