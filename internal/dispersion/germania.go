package dispersion

import (
	"fmt"
	"math"
)

func Germania() Sellmeier {
	return Sellmeier{
		Name: "vitreous GeO2 (Fleming)",
		Terms: []Term{
			{B: 0.80686642, C: 0.006876193},
			{B: 0.71815848, C: 0.010243149},
			{B: 0.85416831, C: 97.933534},
		},
	}
}

func MixSellmeier(silica, dopant Sellmeier, moleFrac float64) (Sellmeier, error) {
	if moleFrac < 0 || moleFrac > 0.4 {
		return Sellmeier{}, fmt.Errorf("dispersion: GeO2 mole fraction %v out of [0, 0.4]", moleFrac)
	}
	if len(silica.Terms) != len(dopant.Terms) {
		return Sellmeier{}, fmt.Errorf("dispersion: Sellmeier term count mismatch")
	}
	x := moleFrac
	terms := silica.Terms
	for i := range terms {
		sb := terms[i].B
		sc := terms[i].C
		terms[i].B = sb*(1-x) + dopant.Terms[i].B*x
		terms[i].C = sc*(1-x) + dopant.Terms[i].C*x
	}
	return Sellmeier{Name: "GeO2-silica mix", Terms: terms}, nil
}

func MaterialDispersionOf(s Sellmeier, lambdaUm float64) float64 {
	n2 := s.SecondDerivativeUm(lambdaUm)
	return -(lambdaUm / SpeedOfLightMPS) * n2 * 1e12
}

func MaterialZeroOf(s Sellmeier, loUm, hiUm float64) (float64, error) {
	if loUm <= 0 || hiUm <= loUm {
		return 0, fmt.Errorf("dispersion: bad zero-search window")
	}
	flo := MaterialDispersionOf(s, loUm)
	fhi := MaterialDispersionOf(s, hiUm)
	if flo == 0 {
		return loUm, nil
	}
	if fhi == 0 {
		return hiUm, nil
	}
	if (flo < 0) == (fhi < 0) {
		return 0, fmt.Errorf("dispersion: no material zero in [%v, %v] μm", loUm, hiUm)
	}
	lo, hi := loUm, hiUm
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		fm := MaterialDispersionOf(s, mid)
		if fm == 0 || hi-lo < 1e-6 {
			return mid, nil
		}
		if (fm < 0) == (flo < 0) {
			lo = mid
			flo = fm
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

func IndexContrast(core, clad Sellmeier, lambdaUm float64) float64 {
	n1 := core.IndexUm(lambdaUm)
	n2 := clad.IndexUm(lambdaUm)
	return n1 - n2
}

func DeltaOf(core, clad Sellmeier, lambdaUm float64) (float64, error) {
	n1 := core.IndexUm(lambdaUm)
	n2 := clad.IndexUm(lambdaUm)
	if n1 <= n2 {
		return 0, fmt.Errorf("dispersion: core index %v not above clad %v at %v μm", n1, n2, lambdaUm)
	}
	return (n1 - n2) / n1, nil
}

func GermaniaRaisesIndex(lambdaUm, moleFrac float64) (float64, error) {
	mix, err := MixSellmeier(Silica(), Germania(), moleFrac)
	if err != nil {
		return 0, err
	}
	return mix.IndexUm(lambdaUm) - Silica().IndexUm(lambdaUm), nil
}

func GermaniaMoleForDelta(lambdaUm, targetDelta float64) (float64, error) {
	if targetDelta <= 0 || targetDelta > 0.02 {
		return 0, fmt.Errorf("dispersion: target Δ %v out of (0, 0.02]", targetDelta)
	}
	lo, hi := 0.0, 0.4
	clad := Silica()
	for i := 0; i < 50; i++ {
		mid := 0.5 * (lo + hi)
		core, err := MixSellmeier(clad, Germania(), mid)
		if err != nil {
			return 0, err
		}
		d, err := DeltaOf(core, clad, lambdaUm)
		if err != nil {
			lo = mid
			continue
		}
		if math.Abs(d-targetDelta) < 1e-6 {
			return mid, nil
		}
		if d < targetDelta {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}
