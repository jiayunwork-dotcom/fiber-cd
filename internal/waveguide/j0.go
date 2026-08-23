package waveguide

// J0 用幂级数计算第一类零阶贝塞尔函数：
//
//	J0(x) = Σ_{m=0}^∞ (−1)^m (x/2)^{2m} / (m!)²
//
// 级数在 |x| ≤ 10 时收敛极快，本包只需该区间的值（截止根 2.405
// 就在其中）。逐项累加到项不再改变求和为止。
func J0(x float64) float64 {
	if x < 0 {
		x = -x
	}
	x2 := x * x / 4
	term := 1.0 // m = 0 项
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

// J0FirstRoot 用二分法求 J0 的首个正根。精确值为 2.4048255577…，
// 圆阶跃光纤的截止常数 V=2.405（CutoffV）就是它的工程舍入。
func J0FirstRoot() float64 {
	lo, hi := 2.0, 3.0 // J0(2) > 0，J0(3) < 0
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

// CutoffVFromJ0 返回 J0 首根，与常量 CutoffV=2.405 互为印证。
func CutoffVFromJ0() float64 {
	return J0FirstRoot()
}
