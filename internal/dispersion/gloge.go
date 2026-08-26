package dispersion

const (
	GlogeConstant = 1.1428
	GlogeSlope    = 0.996
)

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func NormalizedB(v float64) float64 {
	if v <= 0 {
		return 0
	}
	u := GlogeConstant - GlogeSlope/v
	b := u * u
	return clamp(b, 0, 1)
}

func VDBdVSq(v float64) float64 {
	vv := clamp(v, 1.5, 2.5)
	u := GlogeConstant - GlogeSlope/vv

	dbdv := 2 * u * GlogeSlope / (vv * vv)
	d2bdv2 := 2*GlogeSlope*GlogeSlope/(vv*vv*vv*vv) - 4*u*GlogeSlope/(vv*vv*vv)
	d2vb := 2*dbdv + vv*d2bdv2
	return vv * d2vb
}
