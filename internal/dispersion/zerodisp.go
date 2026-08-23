package dispersion

import (
	"errors"
	"fmt"

	"fiber-cd/internal/model"
)

// ErrNoZeroDispersion 表示给定搜索区间内总色散 D_tot 没有跨零
// （两端点同号），无法定位零色散波长。
var ErrNoZeroDispersion = errors.New("no zero-dispersion crossing in range")

// ZeroDispersionWavelength 在 [loNm, hiNm]（nm）内用二分法求总色散
// D_tot(λ) = D_mat(λ) + D_wg(λ) 的零点，返回零色散波长（nm）。
//
// 前提：D_tot 在区间两端点异号。两端同号时返回 ErrNoZeroDispersion，
// 调用方（sweep）据此把结果标为"未发现"而不是编造一个数值。
func ZeroDispersionWavelength(cfg model.Config, loNm, hiNm float64) (float64, error) {
	if err := model.Validate(cfg); err != nil {
		return 0, err
	}
	if loNm <= 0 || hiNm <= loNm {
		return 0, fmt.Errorf("invalid search range [%g, %g] nm", loNm, hiNm)
	}

	na := NumericalAperture(cfg.N1, cfg.N2)
	a := cfg.CoreRadiusM()

	f := func(lamNm float64) float64 {
		lamM := lamNm * 1e-9
		v := vAt(a, na, lamM)
		r := dispersionAt(cfg, v, lamM)
		return r.DTotal
	}

	flo := f(loNm)
	fhi := f(hiNm)
	if flo == 0 {
		return loNm, nil
	}
	if fhi == 0 {
		return hiNm, nil
	}
	if (flo < 0) == (fhi < 0) {
		return 0, flattenZeroErr(ErrNoZeroDispersion)
	}

	lo, hi := loNm, hiNm
	for i := 0; i < 100; i++ {
		mid := (lo + hi) / 2
		fmid := f(mid)
		if fmid == 0 || (hi-lo) < 1e-3 {
			return mid, nil
		}
		if (fmid < 0) == (flo < 0) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2, nil
}
