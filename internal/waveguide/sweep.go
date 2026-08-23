package waveguide

import (
	"errors"
	"fmt"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
)

// SweepStep 是波长扫描表的一行：一个波长档位下的 V、模式状态、色散
// 与色散斜率。
type SweepStep struct {
	WavelengthNm float64
	V            float64
	Mode         ModeStatus
	Disp         dispersion.Result
	// Slope 总色散斜率 S_tot（ps/(nm²·km)），中心差分求得。
	Slope float64
}

// SweepResult 汇总一次波长扫描的全部产物：逐档数据、单模边界与
// 零色散波长（若区间内跨零）。
type SweepResult struct {
	// Steps 逐档扫描数据，波长严格递增。
	Steps []SweepStep
	// StartNm / StopNm 实际扫描区间（nm）。
	StartNm, StopNm float64
	// CutoffNm 单模截止波长（nm，解析值）。CutoffFound 表示它落在
	// 扫描区间内。
	CutoffNm    float64
	CutoffFound bool
	// ZeroDNm 零色散波长（nm）。ZeroDFound 表示区间内存在跨零。
	ZeroDNm    float64
	ZeroDFound bool
}

// Sweep 在 [startNm, stopNm] 上取 steps 档（含两端点）逐波长核算
// V、模式状态与三项色散。波长递增时 V 递减：低于截止的档必须标
// 多模，跨过 V=2.405 之后才允许标单模。
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
		res.Steps = append(res.Steps, SweepStep{
			WavelengthNm: lamNm,
			V:            v,
			Mode:         Classify(v),
			Disp:         d,
			Slope:        slope,
		})
	}

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

// MinAbsDispersion 返回扫描表里 |D_tot| 最小的档位及其波长，用于
// 定位色散最平坦的工作点。
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

// SingleModeWindow 返回扫描区间内保持单模的波长范围 [lo, hi]（nm）。
// 波长递增时 V 递减，单模区一定是连续的一段；扫描全在多模时
// found=false。
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

// absFloat 返回绝对值，避免引入 math 依赖噪音。
func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
