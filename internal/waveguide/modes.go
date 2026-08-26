package waveguide

import "fmt"

func TotalModeCountEstimate(v float64) float64 {
	return v * v / 2
}

func EffectiveModeCount(v float64) float64 {
	if Classify(v) == SingleMode {
		return 1
	}
	return TotalModeCountEstimate(v)
}

func ModeCountDescription(v float64) string {
	if Classify(v) == SingleMode {
		return "1 mode (LP01, single-mode)"
	}
	return "multi-mode (V^2/2 ≈ " + formatCount(TotalModeCountEstimate(v)) + " modes)"
}

func formatCount(n float64) string {
	return fmt.Sprintf("%.0f", n)
}
