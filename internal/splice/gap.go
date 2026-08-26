package splice

import (
	"fmt"
	"math"

	"fiber-cd/internal/waveguide"
)

func GapLossDB(mfdUm, gapUm, wavelengthNm float64) (float64, error) {
	if mfdUm <= 0 || wavelengthNm <= 0 {
		return 0, fmt.Errorf("splice: MFD and wavelength must be > 0")
	}
	if gapUm < 0 {
		return 0, fmt.Errorf("splice: gap must be >= 0")
	}
	w := mfdUm / 2
	zr := math.Pi * w * w / (wavelengthNm / 1000)
	if zr <= 0 {
		return 0, fmt.Errorf("splice: Rayleigh range vanished")
	}
	frac := 1 / (1 + (gapUm/(2*zr))*(gapUm/(2*zr)))
	if frac <= 0 {
		return 0, fmt.Errorf("splice: gap coupling vanished")
	}
	return -10 * math.Log10(frac), nil
}

func CombinedWithGap(overlap, offset, angle, gap float64) (float64, error) {
	for _, v := range []float64{overlap, offset, angle, gap} {
		if v < 0 {
			return 0, fmt.Errorf("splice: loss terms must be >= 0")
		}
	}
	return overlap + offset + angle + gap, nil
}

func EvaluateWithDefinitions(p Pair, offsetUm float64, petermann bool) (Loss, error) {
	a, err := waveguide.Analyze(p.A)
	if err != nil {
		return Loss{}, err
	}
	b, err := waveguide.Analyze(p.B)
	if err != nil {
		return Loss{}, err
	}
	mfd1 := waveguide.ChooseMFDUm(p.A.RadiusUm(), a.V, petermann)
	mfd2 := waveguide.ChooseMFDUm(p.B.RadiusUm(), b.V, petermann)
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

func MFDDefinitionDeltaDB(p Pair, offsetUm float64) (float64, error) {
	marc, err := EvaluateWithDefinitions(p, offsetUm, false)
	if err != nil {
		return 0, err
	}
	pete, err := EvaluateWithDefinitions(p, offsetUm, true)
	if err != nil {
		return 0, err
	}
	return pete.Total - marc.Total, nil
}

func FresnelPairDB(n1, n2 float64) (float64, error) {
	if n1 <= 0 || n2 <= 0 {
		return 0, fmt.Errorf("splice: indices must be > 0")
	}
	r := (n1 - n2) / (n1 + n2)
	refl := r * r
	if refl >= 1 {
		return 0, fmt.Errorf("splice: total reflection")
	}
	return -10 * math.Log10(1-refl), nil
}
