package waveguide

func J0(x float64) float64 {
	if x < 0 {
		x = -x
	}
	x2 := x * x / 4
	term := 1.0
	sum := term
	for m := 1; m < 40; m++ {
		term *= -x2 / float64(m*m)
		prev := sum
		sum += term
		if sum == prev {
			break
		}
	}
	return sum
}

func J0FirstRoot() float64 {
	lo, hi := 2.0, 3.0
	for i := 0; i < 80; i++ {
		mid := (lo + hi) / 2
		if J0(mid) > 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func CutoffVFromJ0() float64 {
	return J0FirstRoot()
}
