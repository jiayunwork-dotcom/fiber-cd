package waveguide

import (
	"fmt"
	"math"
)

const (
	PetermannA = 0.616
	PetermannB = 1.660
	PetermannC = 0.987
)

func petermannRatio(v float64) float64 {
	vv := math.Max(v, 1.0)
	return PetermannA + PetermannB/math.Pow(vv, 1.5) + PetermannC/math.Pow(vv, 6)
}

func marcuseRatio(v float64) float64 {
	vv := math.Max(v, 1.0)
	return MarcuseConstant + MarcuseMidTerm/math.Pow(vv, 1.5) + MarcuseTailTerm/math.Pow(vv, 6)
}

func PetermannSpotUm(coreRadiusUm, v float64) float64 {
	return coreRadiusUm * petermannRatio(v)
}

func PetermannMFDUm(coreRadiusUm, v float64) float64 {
	return 2 * PetermannSpotUm(coreRadiusUm, v)
}

func MFDRatioPetermannToMarcuse(v float64) float64 {
	m := marcuseRatio(v)
	if m <= 0 {
		return 1
	}
	return petermannRatio(v) / m
}

func SpotDisagreementUm(coreRadiusUm, v float64) float64 {
	return math.Abs(PetermannSpotUm(coreRadiusUm, v) - SpotRadiusUm(coreRadiusUm, v))
}

func ChooseMFDUm(coreRadiusUm, v float64, petermann bool) float64 {
	if petermann {
		return PetermannMFDUm(coreRadiusUm, v)
	}
	return ModeFieldDiameterUm(coreRadiusUm, v)
}

func EffectiveAreaFromMFD(mfdUm float64) (float64, error) {
	if mfdUm <= 0 {
		return 0, fmt.Errorf("waveguide: MFD must be > 0")
	}
	w := mfdUm / 2
	return math.Pi * w * w, nil
}

func NonlinearCoefficient(mfdUm, n2 float64, wavelengthNm float64) (float64, error) {
	aeff, err := EffectiveAreaFromMFD(mfdUm)
	if err != nil {
		return 0, err
	}
	if n2 <= 0 || wavelengthNm <= 0 {
		return 0, fmt.Errorf("waveguide: n2 and wavelength must be > 0")
	}
	aeffM2 := aeff * 1e-12
	lam := wavelengthNm * 1e-9
	return 2 * math.Pi * n2 / (lam * aeffM2), nil
}
