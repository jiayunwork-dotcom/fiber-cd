package dispersion

import (
	"errors"
	"fmt"

	"fiber-cd/internal/model"
)

var ErrNoZeroDispersion = errors.New("no zero-dispersion crossing in range")

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
		return 0, ErrNoZeroDispersion
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
