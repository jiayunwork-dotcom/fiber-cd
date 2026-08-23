package model

// 单位换算辅助。所有核算统一先折到 SI 基准（米），再按输出需要还原。
// 折算点只有这里，别处不得再写裸的 1e-6/1e-9 魔法数。

const (
	// MetersPerUm 1 μm 对应的米数。
	MetersPerUm = 1e-6
	// MetersPerNm 1 nm 对应的米数。
	MetersPerNm = 1e-9
	// NmPerUm 1 μm 对应的 nm 数。
	NmPerUm = 1000.0
)

// NmToM 纳米转米。
func NmToM(nm float64) float64 { return nm * MetersPerNm }

// UmToM 微米转米。
func UmToM(um float64) float64 { return um * MetersPerUm }

// MToNm 米转纳米。
func MToNm(m float64) float64 { return m / MetersPerNm }

// MToUm 米转微米。
func MToUm(m float64) float64 { return m / MetersPerUm }

// UmToNm 微米转纳米。
func UmToNm(um float64) float64 { return um * NmPerUm }

// NmToUm 纳米转微米。
func NmToUm(nm float64) float64 { return nm / NmPerUm }

// PercentToFraction 百分数（如 0.25 表示 0.25%）转小数比例。
func PercentToFraction(percent float64) float64 {
	return percent / 100
}
