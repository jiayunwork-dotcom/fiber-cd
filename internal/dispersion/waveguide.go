package dispersion

// WaveguideDispersion 计算阶跃近似波导色散：
//
//	D_wg = −(n1·Δ/(λ·c))·V·d²(Vb)/dV²
//
// 这是光纤教材里钉死的文献形式（Gloge 1971 的 b(V) 配合上式）。
// 入参全部为 SI 单位：n1、Δ 无量纲，lambdaM 取米；返回 ps/(nm·km)。
//
// 符号约定：阶跃光纤中 D_wg 恒为负，与材料项在零色散波长附近对消，
// 使总色散跨零翻转。
func WaveguideDispersion(n1, delta, lambdaM, v float64) float64 {
	g := VDBdVSq(v)
	return -(n1 * delta / (lambdaM * SpeedOfLightMPS)) * g * 1e6
}
