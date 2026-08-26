package dispersion

func (s Sellmeier) FirstDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		denom := lambdaUm*lambdaUm - t.C
		sum += t.B * t.C / (denom * denom)
	}
	return -lambdaUm / n * sum
}

func (s Sellmeier) SecondDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	n1 := s.FirstDerivativeUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		denom := l2 - t.C
		cube := denom * denom * denom
		sum += t.C * t.B * (3*l2 + t.C) / cube
	}
	return (sum - n1*n1) / n
}

func (s Sellmeier) ThirdDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	n1 := s.FirstDerivativeUm(lambdaUm)
	n2 := s.SecondDerivativeUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		denom := l2 - t.C
		pow4 := denom * denom * denom * denom
		sum += t.C * t.B * (l2 + t.C) / pow4
	}
	return (-24*lambdaUm*sum - 4*n1*n2 - 2*n2*n2) / (2 * n)
}
