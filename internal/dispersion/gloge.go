package dispersion

// Gloge 1971 给出了阶跃光纤基模归一化传播常数 b(V) 的著名近似：
//
//	b(V) ≈ (1.1428 − 0.996/V)²
//
// 严格适用区间约 1.5 < V < 2.5。域外 b 会越过物理边界 [0, 1]，
// 这里对 b 做钳制；导数量 V·d²(Vb)/dV² 在域外饱和于边界值，
// 保证波导色散处处有限。
const (
	// GlogeConstant 是 Gloge 近似的常数项 1.1428。
	GlogeConstant = 1.1428
	// GlogeSlope 是 Gloge 近似的斜率项 0.996。
	GlogeSlope = 0.996
)

// clamp 把 x 钳制到 [lo, hi]。
func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// NormalizedB 返回归一化传播常数 b(V)，取值被钳制在 [0, 1]。
func NormalizedB(v float64) float64 {
	if v <= 0 {
		return 0
	}
	u := GlogeConstant - GlogeSlope/v
	b := u * u
	return clamp(b, 0, 1)
}

// VDBdVSq 返回波导色散公式里的形状因子 g(V) = V·d²(Vb)/dV²，
// 由 Gloge 近似解析求导得出。V 在 [1.5, 2.5] 外时取边界值饱和。
func VDBdVSq(v float64) float64 {
	vv := clamp(v, 1.5, 2.5)
	u := GlogeConstant - GlogeSlope/vv

	// d b/dV = 2u·(0.996/V²)
	dbdv := 2 * u * GlogeSlope / (vv * vv)
	// d²b/dV² = 2·0.996²/V⁴ − 4u·0.996/V³
	d2bdv2 := 2*GlogeSlope*GlogeSlope/(vv*vv*vv*vv) - 4*u*GlogeSlope/(vv*vv*vv)
	// d²(Vb)/dV² = 2·db/dV + V·d²b/dV²
	d2vb := 2*dbdv + vv*d2bdv2
	return vv * d2vb
}
