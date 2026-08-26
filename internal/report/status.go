package report

import "fiber-cd/internal/waveguide"

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
