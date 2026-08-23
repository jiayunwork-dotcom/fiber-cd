package waveguide

import "math"

// VNumber 计算圆阶跃光纤的归一化频率 V = 2π·a·NA/λ。
//
// 入参一律为 SI 单位（a、λ 取米）。a 是纤芯半径，不是直径：把芯径
// 直接塞进来会让 V 翻倍、截止波长差一倍、单模判定全错。
func VNumber(coreRadiusM, na, wavelengthM float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / wavelengthM
}
