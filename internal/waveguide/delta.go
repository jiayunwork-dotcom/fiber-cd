package waveguide

// RelativeIndexDelta 计算相对折射率差 Δ = (n1 − n2) / n1。
//
// Δ 同时驱动波导色散公式（n1·Δ 项）与折射率剖面描述，是连接导波量
// 与色散量的桥梁。n1 恒大于 n2 时结果落在 (0, 1)。
func RelativeIndexDelta(n1, n2 float64) float64 {
	return (n1 - n2) / n1
}
