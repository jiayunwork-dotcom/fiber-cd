// Package waveguide 是圆阶跃光纤的导波内核：数值孔径 NA、相对折射率差
// Δ、归一化频率 V、单模截止边界（V < 2.405）与截止波长。
//
// 关键约定：用户输入的是芯径（直径），V 数公式 2πa·NA/λ 里的 a 是
// 纤芯半径，统一由 model.Config.RadiusUm() 从芯径折算。
package waveguide

import "math"

// NumericalAperture 计算圆阶跃光纤的数值孔径 NA = √(n1² − n2²)。
//
// 调用方必须保证 n2 < n1（model.Validate 负责）；这里不做二次校验，
// 输入不合法时返回 NaN 属预期行为。
func NumericalAperture(n1, n2 float64) float64 {
	return math.Sqrt(n1*n1 - n2*n2)
}
