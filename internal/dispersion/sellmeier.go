package dispersion

import "math"

type Term struct {
	B float64
	C float64
}

type Sellmeier struct {
	Name  string
	Terms []Term
}

func Silica() Sellmeier {
	return Sellmeier{
		Name: "fused silica (Malitson 1965)",
		Terms: []Term{
			{B: 0.6961663, C: 0.00467914826},
			{B: 0.4079426, C: 0.0135120631},
			{B: 0.8974794, C: 97.9340025},
		},
	}
}

func (s Sellmeier) IndexUm(lambdaUm float64) float64 {
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		sum += t.B * l2 / (l2 - t.C)
	}
	return math.Sqrt(1 + sum)
}
