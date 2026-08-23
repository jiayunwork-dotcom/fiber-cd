package dispersion

import (
	"math"

	"fiber-cd/internal/model"
)

// NumericalAperture 计算 NA = √(n1²−n2²)。
//
// 与 waveguide.NumericalAperture 等价；本包保留一份本地拷贝以维持
// 依赖方向（waveguide → dispersion），避免循环导入。
func NumericalAperture(n1, n2 float64) float64 {
	return math.Sqrt(n1*n1 - n2*n2)
}

// RelativeDelta 计算 Δ = (n1−n2)/n1，供波导色散公式使用。
// 同为 waveguide.RelativeIndexDelta 的本地拷贝。
func RelativeDelta(n1, n2 float64) float64 {
	return (n1 - n2) / n1
}

// vAt 按 V = 2π·a·NA/λ 计算归一化频率，a、λ 取米。a 是半径。
func vAt(coreRadiusM, na, lambdaM float64) float64 {
	return 2 * math.Pi * coreRadiusM * na / lambdaM
}

// dispersionAt 在任意波长（米）下合成三项色散。
func dispersionAt(cfg model.Config, v, lambdaM float64) Result {
	lambdaNm := lambdaM * 1e9
	lambdaUm := lambdaM * 1e6
	dMat := MaterialDispersionUm(lambdaUm)
	dWg := WaveguideDispersion(cfg.N1, RelativeDelta(cfg.N1, cfg.N2), lambdaM, v)
	return Result{
		LambdaNm: lambdaNm,
		V:        v,
		DMat:     dMat,
		DWg:      dWg,
		DTotal:   dMat + dWg,
	}
}

// Compose 对给定配置，在指定的波长（米）与归一化频率 V 下合成色散
// 结果。V 由调用方传入（通常来自 waveguide.VNumber）。
func Compose(cfg model.Config, v, lambdaM float64) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	return dispersionAt(cfg, v, lambdaM), nil
}

// TotalDispersionAt 按配置自身的工作波长完成一次完整色散核算。
// 这是 mode 子命令用的单点入口。
func TotalDispersionAt(cfg model.Config) (Result, error) {
	if err := model.Validate(cfg); err != nil {
		return Result{}, err
	}
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	return dispersionAt(cfg, v, cfg.WavelengthM()), nil
}
