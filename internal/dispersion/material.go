package dispersion

// SpeedOfLightMPS 是真空光速，单位 m/s，取 CODATA 精确值
// 299792458，不在任何换算处四舍五入。
const SpeedOfLightMPS = 299792458.0

// MaterialDispersionUm 计算材料色散
//
//	D_mat = −(λ/c)·d²n/dλ²
//
// 返回单位 ps/(nm·km)，λ 以 μm 输入，n(λ) 来自本包钉死的
// Malitson 1965 石英 Sellmeier 系数。单位换算：SI 结果
// (s/m²) 乘以 10⁶ 即为 ps/(nm·km)。
func MaterialDispersionUm(lambdaUm float64) float64 {
	n2 := Silica().SecondDerivativeUm(lambdaUm)
	return -(lambdaUm / SpeedOfLightMPS) * n2 * 1e12
}

// MaterialDispersionSlopeUm 计算材料色散斜率 S_mat = dD_mat/dλ，
// 单位 ps/(nm²·km)，λ 以 μm 输入。由 D_mat 的闭合形式直接求导：
//
//	S_mat = −(1/c)·(n'' + λ·n''')·10⁹
func MaterialDispersionSlopeUm(lambdaUm float64) float64 {
	silica := Silica()
	n2 := silica.SecondDerivativeUm(lambdaUm)
	n3 := silica.ThirdDerivativeUm(lambdaUm)
	return -(1 / SpeedOfLightMPS) * (n2 + lambdaUm*n3) * 1e9
}
