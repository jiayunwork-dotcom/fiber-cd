package report

import "fiber-cd/internal/waveguide"

// ModeStatusLabel 把模式状态转成面向用户的英文短语，供报告与扫描
// 表使用，保持全仓措辞一致。
func ModeStatusLabel(m waveguide.ModeStatus) string {
	switch m {
	case waveguide.SingleMode:
		return "single-mode"
	case waveguide.MultiMode:
		return "multi-mode"
	default:
		return "unknown"
	}
}
