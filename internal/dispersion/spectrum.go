package dispersion

import (
	"fmt"

	"fiber-cd/internal/model"
)

// SpectrumPoint 是某一波长下材料与导波量的完整光谱点：折射率、
// 群折射率、三项色散与总色散斜率。spec / probe 报告用它做聚合输出。
type SpectrumPoint struct {
	LambdaNm   float64
	LambdaUm   float64
	Index      float64
	GroupIndex float64
	DMat       float64
	DWg        float64
	DTotal     float64
	STotal     float64
}

// SpectrumAt 在指定波长（nm）下合成一个完整光谱点：V 按配置几何
// 与波长即时重算，色散按该波长求值。
func SpectrumAt(cfg model.Config, lambdaNm float64) (SpectrumPoint, error) {
	if err := model.Validate(cfg); err != nil {
		return SpectrumPoint{}, err
	}
	if lambdaNm <= 0 {
		return SpectrumPoint{}, fmt.Errorf("wavelength must be positive, got %g nm", lambdaNm)
	}
	lambdaM := lambdaNm * 1e-9
	lambdaUm := lambdaNm / 1000

	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, lambdaM)
	d := dispersionAt(cfg, v, lambdaM)
	s, err := TotalDispersionSlope(cfg, lambdaNm)
	if err != nil {
		return SpectrumPoint{}, err
	}

	silica := Silica()
	return SpectrumPoint{
		LambdaNm:   lambdaNm,
		LambdaUm:   lambdaUm,
		Index:      silica.IndexUm(lambdaUm),
		GroupIndex: GroupIndex(lambdaUm),
		DMat:       d.DMat,
		DWg:        d.DWg,
		DTotal:     d.DTotal,
		STotal:     s,
	}, nil
}
