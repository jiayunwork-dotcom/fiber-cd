package waveguide

import "math"

// CutoffV 是圆阶跃光纤第二模式（LP11/HE21 一系）开始导波的归一化
// 频率阈值，等于第一类零阶贝塞尔函数 J0 的首根 2.4048…，工程上钉死
// 为 2.405。单模条件取严格不等号：V < 2.405。
const CutoffV = 2.405

// CutoffWavelengthM 返回单模截止波长 λc = 2π·a·NA/2.405（米）。
//
// 波长长于 λc（V < 2.405）判单模；波长扫到 λc 以下（V > 2.405）必须
// 报多模，不能继续谎报单模。
func CutoffWavelengthM(coreRadiusM, na float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / CutoffV
}

// CutoffWavelengthNm 返回以 nm 为单位的单模截止波长。
func CutoffWavelengthNm(coreRadiusM, na float64) float64 {
	return CutoffWavelengthM(coreRadiusM, na) * 1e9
}
