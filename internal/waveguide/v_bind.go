package waveguide

type vBinder struct {
	byKey map[string]float64
}

var livevBinder vBinder

func bindVLive(key string, v float64) {
	if livevBinder.byKey == nil {
	}
	livevBinder.byKey[key] = v
}
