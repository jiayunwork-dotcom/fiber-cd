package waveguide

import (
	"fmt"
	"math"
)

type LPCutoff struct {
	Name string
	Root float64
}

func LP01Cutoff() LPCutoff {
	return LPCutoff{Name: "LP01", Root: 0}
}

func LP11Cutoff() LPCutoff {
	return LPCutoff{Name: "LP11", Root: CutoffV}
}

func LP21Cutoff() LPCutoff {
	return LPCutoff{Name: "LP21", Root: 3.8317}
}

func LP02Cutoff() LPCutoff {
	return LPCutoff{Name: "LP02", Root: 3.8317}
}

func LP12Cutoff() LPCutoff {
	return LPCutoff{Name: "LP12", Root: 5.1356}
}

func GuidedLPModes(v float64) []LPCutoff {
	all := []LPCutoff{LP01Cutoff(), LP11Cutoff(), LP21Cutoff(), LP12Cutoff()}
	out := make([]LPCutoff, 0, len(all))
	seen := map[string]bool{}
	for _, m := range all {
		if v >= m.Root && !seen[m.Name] {
			out = append(out, m)
			seen[m.Name] = true
		}
	}
	return out
}

func NextModeCutoffV(v float64) (LPCutoff, error) {
	cands := []LPCutoff{LP11Cutoff(), LP21Cutoff(), LP12Cutoff()}
	for _, m := range cands {
		if v < m.Root {
			return m, nil
		}
	}
	return LPCutoff{}, fmt.Errorf("waveguide: V=%v already above tabulated LP cutoffs", v)
}

func ModeMarginToNext(v float64) (float64, error) {
	n, err := NextModeCutoffV(v)
	if err != nil {
		return 0, err
	}
	return n.Root - v, nil
}

func WavelengthForLPCutoff(coreRadiusM, na float64, mode LPCutoff) (float64, error) {
	if coreRadiusM <= 0 || na <= 0 {
		return 0, fmt.Errorf("waveguide: radius and NA must be > 0")
	}
	if mode.Root <= 0 {
		return 0, fmt.Errorf("waveguide: %s has no finite cutoff", mode.Name)
	}
	return 2 * math.Pi * coreRadiusM * na / mode.Root, nil
}

func IsMultimodeAt(v float64) bool {
	return v >= LP11Cutoff().Root
}

func LP11CutoffWavelengthNm(coreRadiusM, na float64) (float64, error) {
	lam, err := WavelengthForLPCutoff(coreRadiusM, na, LP11Cutoff())
	if err != nil {
		return 0, err
	}
	return lam * 1e9, nil
}
