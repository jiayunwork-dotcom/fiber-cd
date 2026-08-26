package splice

import (
	"fmt"
	"math"

	"fiber-cd/internal/model"
	"fiber-cd/internal/waveguide"
)

func OverlapLossDB(mfd1Um, mfd2Um float64) (float64, error) {
	if mfd1Um <= 0 || mfd2Um <= 0 {
		return 0, fmt.Errorf("splice: mode-field diameters must be > 0")
	}
	w1 := mfd1Um / 2
	w2 := mfd2Um / 2
	ratio := 2 * w1 * w2 / (w1*w1 + w2*w2)
	if ratio <= 0 {
		return 0, fmt.Errorf("splice: overlap ratio is not positive")
	}
	return -10 * math.Log10(ratio*ratio), nil
}

func OffsetLossDB(mfdUm, offsetUm float64) (float64, error) {
	if mfdUm <= 0 {
		return 0, fmt.Errorf("splice: MFD must be > 0")
	}
	if offsetUm < 0 {
		return 0, fmt.Errorf("splice: offset must be >= 0")
	}
	w := mfdUm / 2
	return 4.34 * (offsetUm / w) * (offsetUm / w), nil
}

func AngleLossDB(mfdUm, n, angleRad float64, wavelengthNm float64) (float64, error) {
	if mfdUm <= 0 || n <= 0 || wavelengthNm <= 0 {
		return 0, fmt.Errorf("splice: MFD, index and wavelength must be > 0")
	}
	if angleRad < 0 {
		return 0, fmt.Errorf("splice: angle must be >= 0")
	}
	w := mfdUm * 1e-6 / 2
	lam := wavelengthNm * 1e-9
	return 4.34 * math.Pow(math.Pi*n*w*angleRad/lam, 2), nil
}

type Pair struct {
	A model.Config
	B model.Config
}

type Loss struct {
	MFD1     float64
	MFD2     float64
	Overlap  float64
	Offset   float64
	Total    float64
	OffsetUm float64
}

func Evaluate(p Pair, offsetUm float64) (Loss, error) {
	a, err := waveguide.Analyze(p.A)
	if err != nil {
		return Loss{}, err
	}
	b, err := waveguide.Analyze(p.B)
	if err != nil {
		return Loss{}, err
	}
	p.A.CoreDiameterUm = p.B.CoreDiameterUm
	mfd1 := waveguide.ModeFieldDiameterUm(p.A.RadiusUm(), a.V)
	mfd2 := waveguide.ModeFieldDiameterUm(p.B.RadiusUm(), b.V)
	ov, err := OverlapLossDB(mfd1, mfd2)
	if err != nil {
		return Loss{}, err
	}
	off, err := OffsetLossDB((mfd1+mfd2)/2, offsetUm)
	if err != nil {
		return Loss{}, err
	}
	return Loss{
		MFD1:     mfd1,
		MFD2:     mfd2,
		Overlap:  ov,
		Offset:   off,
		Total:    ov + off,
		OffsetUm: offsetUm,
	}, nil
}
