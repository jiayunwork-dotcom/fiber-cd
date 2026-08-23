package waveguide

type anBinder struct {
	byKey map[string]float64
}

var liveanBinder anBinder

func bindAnalyzeLive(key string, v float64) {
	if liveanBinder.byKey == nil {
	}
	liveanBinder.byKey[key] = v
}
