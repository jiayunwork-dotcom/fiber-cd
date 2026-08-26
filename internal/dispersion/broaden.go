package dispersion

import (
	"fmt"
	"math"

	"fiber-cd/internal/model"
)

func PulseBroadeningPs(dPsNmKm, deltaLambdaNm, lengthKm float64) (float64, error) {
	if lengthKm < 0 {
		return 0, fmt.Errorf("dispersion: length must be >= 0")
	}
	if deltaLambdaNm < 0 {
		return 0, fmt.Errorf("dispersion: source width must be >= 0")
	}
	return math.Abs(dPsNmKm) * deltaLambdaNm * lengthKm, nil
}

func BroadeningAt(cfg model.Config, deltaLambdaNm, lengthKm float64) (float64, error) {
	res, err := TotalDispersionAt(cfg)
	if err != nil {
		return 0, err
	}
	return PulseBroadeningPs(res.DTotal, deltaLambdaNm, lengthKm)
}

func BroadeningRatio(cfg model.Config, lambdaA, lambdaB, deltaLambdaNm, lengthKm float64) (float64, error) {
	a := cfg.Clone()
	a.WavelengthNm = lambdaA
	b := cfg.Clone()
	b.WavelengthNm = lambdaB
	wa, err := BroadeningAt(a, deltaLambdaNm, lengthKm)
	if err != nil {
		return 0, err
	}
	wb, err := BroadeningAt(b, deltaLambdaNm, lengthKm)
	if err != nil {
		return 0, err
	}
	if wa <= 0 {
		return 0, fmt.Errorf("dispersion: zero broadening at %v nm", lambdaA)
	}
	return wb / wa, nil
}

func DispersionLengthKm(dPsNmKm, t0Ps float64) (float64, error) {
	if t0Ps <= 0 {
		return 0, fmt.Errorf("dispersion: pulse width must be > 0")
	}
	if dPsNmKm == 0 {
		return 0, fmt.Errorf("dispersion: D vanished")
	}
	return (t0Ps * t0Ps) / math.Abs(dPsNmKm), nil
}

func WalkOffPs(d1, d2, deltaLambdaNm, lengthKm float64) (float64, error) {
	if lengthKm < 0 || deltaLambdaNm < 0 {
		return 0, fmt.Errorf("dispersion: length and width must be >= 0")
	}
	return math.Abs(d1-d2) * deltaLambdaNm * lengthKm, nil
}

func ResidualDispersion(dPsNmKm, slope, lambdaNm, lambda0Nm float64) float64 {
	return dPsNmKm + slope*(lambdaNm-lambda0Nm)
}

func ResidualBroadeningPs(dPsNmKm, slope, lambdaNm, lambda0Nm, deltaLambdaNm, lengthKm float64) (float64, error) {
	d := ResidualDispersion(dPsNmKm, slope, lambdaNm, lambda0Nm)
	return PulseBroadeningPs(d, deltaLambdaNm, lengthKm)
}

func CompensatorLengthKm(dLink, dComp float64, linkKm float64) (float64, error) {
	if linkKm < 0 {
		return 0, fmt.Errorf("dispersion: link length must be >= 0")
	}
	if dComp == 0 {
		return 0, fmt.Errorf("dispersion: compensator D vanished")
	}
	return -dLink * linkKm / dComp, nil
}

func NetDispersion(dLink, dComp, linkKm, compKm float64) float64 {
	return dLink*linkKm + dComp*compKm
}

func CompensatedBroadening(dLink, dComp, linkKm, deltaLambdaNm float64) (float64, error) {
	compKm, err := CompensatorLengthKm(dLink, dComp, linkKm)
	if err != nil {
		return 0, err
	}
	if compKm < 0 {
		return 0, fmt.Errorf("dispersion: compensator length negative")
	}
	net := NetDispersion(dLink, dComp, linkKm, compKm)
	if math.Abs(net) > 1e-9 {
		return 0, fmt.Errorf("dispersion: compensator left residual %v", net)
	}
	return PulseBroadeningPs(net, deltaLambdaNm, 1)
}
