package dispersion

const SpeedOfLightMPS = 299792458.0

func MaterialDispersionUm(lambdaUm float64) float64 {
	n2 := Silica().SecondDerivativeUm(lambdaUm)
	return -(lambdaUm / SpeedOfLightMPS) * n2 * 1e12
}

func MaterialDispersionSlopeUm(lambdaUm float64) float64 {
	silica := Silica()
	n2 := silica.SecondDerivativeUm(lambdaUm)
	n3 := silica.ThirdDerivativeUm(lambdaUm)
	return -(1 / SpeedOfLightMPS) * (n2 + lambdaUm*n3) * 1e9
}
