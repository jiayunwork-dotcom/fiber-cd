package waveguide

import "math"

// Marcuse 1969 给出了阶跃光纤基模的光斑半径（场半径）近似：
//
//	w0/a = 0.65 + 1.619·V^{−3/2} + 2.879·V^{−6}
//
// 适用于单模区 1.5 < V < 2.5。V 过小会使幂项发散，这里把 V 下限
// 钳到 1.0，保证计算处处有限。
const (
	// MarcuseConstant 是 Marcuse 公式的常数项 0.65。
	MarcuseConstant = 0.65
	// MarcuseMidTerm 是 V^{-3/2} 项的系数 1.619。
	MarcuseMidTerm = 1.619
	// MarcuseTailTerm 是 V^{-6} 项的系数 2.879。
	MarcuseTailTerm = 2.879
)

// SpotRadiusUm 返回基模光斑半径 w0（μm）。coreRadiusUm 为纤芯半径。
func SpotRadiusUm(coreRadiusUm, v float64) float64 {
	vv := math.Max(v, 1.0)
	ratio := MarcuseConstant + MarcuseMidTerm/math.Pow(vv, 1.5) + MarcuseTailTerm/math.Pow(vv, 6)
	return coreRadiusUm * ratio
}

// ModeFieldDiameterUm 返回模场直径 MFD = 2·w0（μm）。
func ModeFieldDiameterUm(coreRadiusUm, v float64) float64 {
	return 2 * SpotRadiusUm(coreRadiusUm, v)
}

// EffectiveAreaUm2 返回基模有效面积 Aeff ≈ π·w0²（μm²）。
func EffectiveAreaUm2(coreRadiusUm, v float64) float64 {
	w0 := SpotRadiusUm(coreRadiusUm, v)
	return math.Pi * w0 * w0
}
