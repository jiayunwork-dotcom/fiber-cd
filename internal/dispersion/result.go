package dispersion

type Result struct {
	LambdaNm float64
	V        float64
	DMat     float64
	DWg      float64
	DTotal   float64
}

func (r Result) DispersionTerms() [3]float64 {
	return [3]float64{r.DMat, r.DWg, r.DTotal}
}

func Terms() [3]string {
	return [3]string{"material", "waveguide", "total"}
}
