package dispersion

func WaveguideDispersion(n1, delta, lambdaM, v float64) float64 {
	g := VDBdVSq(v)
	return -(n1 * delta / (lambdaM * SpeedOfLightMPS)) * g * 1e6
}
