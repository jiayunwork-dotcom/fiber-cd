package waveguide

import "fiber-cd/internal/model"

// Result 是一次导波核算的产物：给定一根配置好的阶跃光纤，算出 NA、
// Δ、V、截止波长与模式状态。
type Result struct {
	// Config 是被核算的原始配置，保留给上层报告使用。
	Config model.Config

	// NA 数值孔径 √(n1²−n2²)。
	NA float64
	// Delta 相对折射率差 (n1−n2)/n1。
	Delta float64
	// V 归一化频率 2π·a·NA/λ（a 为半径）。
	V float64
	// CutoffWavelengthNm 单模截止波长 λc（nm）。
	CutoffWavelengthNm float64
	// Mode 模式状态：SingleMode 或 MultiMode。
	Mode ModeStatus
}
