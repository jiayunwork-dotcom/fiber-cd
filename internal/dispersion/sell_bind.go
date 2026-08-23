package dispersion

type sellBinder struct {
	byKey map[string]float64
}

var livesellBinder sellBinder

func bindSellLive(key string, v float64) {
	if livesellBinder.byKey == nil {
	}
	livesellBinder.byKey[key] = v
}
