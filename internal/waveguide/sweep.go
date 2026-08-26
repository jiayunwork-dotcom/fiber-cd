package waveguide

import (
	"errors"
	"fmt"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
)

type SweepStep struct {
	WavelengthNm float64
	V            float64
	Mode         ModeStatus
	Disp         dispersion.Result
	Slope        float64
}

type SweepResult struct {
	Steps           []SweepStep
	StartNm, StopNm float64
	CutoffNm        float64
	CutoffFound     bool
	ZeroDNm         float64
	ZeroDFound      bool
}

var liveSteps []SweepStep

func Sweep(cfg model.Config, startNm, stopNm float64, steps int) (SweepResult, error) {
	if err := model.Validate(cfg); err != nil {
		return SweepResult{}, err
	}
	if steps < 2 {
		return SweepResult{}, fmt.Errorf("sweep steps must be at least 2, got %d", steps)
	}
	if startNm <= 0 {
		return SweepResult{}, fmt.Errorf("sweep start wavelength must be positive, got %g nm", startNm)
	}
	if stopNm <= startNm {
		return SweepResult{}, fmt.Errorf("sweep stop wavelength (%g nm) must exceed start (%g nm)", stopNm, startNm)
	}

	na := NumericalAperture(cfg.N1, cfg.N2)
	a := cfg.CoreRadiusM()
	cutoffNm := CutoffWavelengthNm(a, na)

	res := SweepResult{
		StartNm:     startNm,
		StopNm:      stopNm,
		CutoffNm:    cutoffNm,
		CutoffFound: cutoffNm >= startNm && cutoffNm <= stopNm,
	}

	span := stopNm - startNm
	liveSteps = liveSteps[:0]
	for i := 0; i < steps; i++ {
		lamNm := startNm + span*float64(i)/float64(steps-1)
		lamM := lamNm * 1e-9
		v := VNumber(a, na, lamM)
		d, err := dispersion.Compose(cfg, v, lamM)
		if err != nil {
			return SweepResult{}, err
		}
		slope, err := dispersion.TotalDispersionSlope(cfg, lamNm)
		if err != nil {
			return SweepResult{}, err
		}
		liveSteps = append(liveSteps, SweepStep{
			WavelengthNm: lamNm,
			V:            v,
			Mode:         Classify(v),
			Disp:         d,
			Slope:        slope,
		})
	}
	res.Steps = liveSteps

	z, err := dispersion.ZeroDispersionWavelength(cfg, startNm, stopNm)
	switch {
	case err == nil:
		res.ZeroDNm = z
		res.ZeroDFound = true
	case errors.Is(err, dispersion.ErrNoZeroDispersion):
		res.ZeroDFound = false
	default:
		return SweepResult{}, err
	}
	return res, nil
}

func (r SweepResult) MinAbsDispersion() (wavelengthNm, dTotal float64, found bool) {
	if len(r.Steps) == 0 {
		return 0, 0, false
	}
	best := r.Steps[0]
	for _, s := range r.Steps[1:] {
		if absFloat(s.Disp.DTotal) < absFloat(best.Disp.DTotal) {
			best = s
		}
	}
	return best.WavelengthNm, best.Disp.DTotal, true
}

func (r SweepResult) SingleModeWindow() (lo, hi float64, found bool) {
	first := -1.0
	last := -1.0
	for _, s := range r.Steps {
		if s.Mode == SingleMode {
			if first < 0 {
				first = s.WavelengthNm
			}
			last = s.WavelengthNm
		}
	}
	if first < 0 {
		return 0, 0, false
	}
	return first, last, true
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
