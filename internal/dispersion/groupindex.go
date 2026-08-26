package dispersion

func GroupIndex(lambdaUm float64) float64 {
	s := Silica()
	return s.IndexUm(lambdaUm) - lambdaUm*s.FirstDerivativeUm(lambdaUm)
}

func GroupDelayPerKm(lambdaUm float64) float64 {
	return GroupIndex(lambdaUm) / SpeedOfLightMPS * 1e9
}

func GroupVelocity(lambdaUm float64) float64 {
	return SpeedOfLightMPS / GroupIndex(lambdaUm)
}
