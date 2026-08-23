package waveguide

import "fmt"

// TotalModeCountEstimate 返回 Gloge 阶跃光纤总模式数近似 M ≈ V²/2
// （双偏振统计，V 数越大越准）。
func TotalModeCountEstimate(v float64) float64 {
	return v * v / 2
}

// EffectiveModeCount 返回可承载的有效模式数：单模为 1（基模），
// 多模用 Gloge 估计值。
func EffectiveModeCount(v float64) float64 {
	if Classify(v) == SingleMode {
		return 1
	}
	return TotalModeCountEstimate(v)
}

// ModeCountDescription 生成一行便于报告的模式数描述。
func ModeCountDescription(v float64) string {
	if Classify(v) == SingleMode {
		return "1 mode (LP01, single-mode)"
	}
	return "multi-mode (V^2/2 ≈ " + formatCount(TotalModeCountEstimate(v)) + " modes)"
}

// formatCount 把模式数四舍五入成整数字符串。
func formatCount(n float64) string {
	return fmt.Sprintf("%.0f", n)
}
