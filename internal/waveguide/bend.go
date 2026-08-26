package waveguide

import (
	"fmt"
	"math"

	"fiber-cd/internal/model"
)

func claddingW(v float64) float64 {
	if v <= 0 {
		return 0
	}
	u := 1.1428 - 0.996/math.Max(v, 0.8)
	if u < 0 {
		u = 0
	}
	b := u * u
	if b > 1 {
		b = 1
	}
	return v * math.Sqrt(math.Max(1-b, 1e-6))
}

func MacrobendLossDBperM(coreRadiusM, v, delta, bendRadiusM float64) (float64, error) {
	if coreRadiusM <= 0 || bendRadiusM <= 0 {
		return 0, fmt.Errorf("waveguide: core and bend radii must be > 0")
	}
	if v <= 0 || delta <= 0 {
		return 0, fmt.Errorf("waveguide: V and Δ must be > 0")
	}
	w := claddingW(v)
	if w <= 0 {
		return 0, fmt.Errorf("waveguide: cladding decay vanished")
	}
	arg := (4.0 / 3.0) * (bendRadiusM / coreRadiusM) * math.Pow(2*delta, 1.5) * w * w * w
	pre := math.Sqrt(math.Pi/(coreRadiusM*bendRadiusM)) / math.Pow(w, 1.5)
	lin := pre * math.Exp(-arg)
	if lin <= 0 || math.IsNaN(lin) || math.IsInf(lin, 0) {
		return 0, fmt.Errorf("waveguide: bend loss not finite")
	}
	return 4.343 * lin, nil
}

func MacrobendLossFor(cfg model.Config, bendRadiusM float64) (float64, error) {
	r, err := Analyze(cfg)
	if err != nil {
		return 0, err
	}
	return MacrobendLossDBperM(cfg.CoreRadiusM(), r.V, r.Delta, bendRadiusM)
}

func CriticalBendRadius(coreRadiusM, v, delta float64) (float64, error) {
	if coreRadiusM <= 0 || v <= 0 || delta <= 0 {
		return 0, fmt.Errorf("waveguide: core, V and Δ must be > 0")
	}
	w := claddingW(v)
	if w <= 0 {
		return 0, fmt.Errorf("waveguide: cladding decay vanished")
	}
	return 1.5 * coreRadiusM / (math.Pow(2*delta, 1.5) * w * w * w), nil
}

func BendLossRatio(cfg model.Config, rLoose, rTight float64) (float64, error) {
	a, err := MacrobendLossFor(cfg, rLoose)
	if err != nil {
		return 0, err
	}
	b, err := MacrobendLossFor(cfg, rTight)
	if err != nil {
		return 0, err
	}
	if a <= 0 {
		return 0, fmt.Errorf("waveguide: loose-bend loss vanished")
	}
	return b / a, nil
}

func NearCutoffBendPenalty(cfg model.Config, bendRadiusM, lambdaNear, lambdaFar float64) (float64, error) {
	near := cfg.Clone()
	near.WavelengthNm = lambdaNear
	far := cfg.Clone()
	far.WavelengthNm = lambdaFar
	ln, err := MacrobendLossFor(near, bendRadiusM)
	if err != nil {
		return 0, err
	}
	lf, err := MacrobendLossFor(far, bendRadiusM)
	if err != nil {
		return 0, err
	}
	if lf <= 0 {
		return 0, fmt.Errorf("waveguide: far-from-cutoff bend loss vanished")
	}
	return ln / lf, nil
}
