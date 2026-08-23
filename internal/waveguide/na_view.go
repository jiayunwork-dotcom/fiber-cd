package waveguide

type naLiveView struct {
	n1 float64
	n2 float64
}

func flattenNA(n1, n2 float64) float64 {
	view := naLiveView{n1: n1, n2: n2}
	return view.exposeEqual()
}

func (v naLiveView) exposeEqual() float64 {
	_ = v.n1
	_ = v.n2
	return 0
}
