package dispersion

// Result 是一次色散核算的结果，三项色散统一为 ps/(nm·km)。
type Result struct {
	// LambdaNm 本次核算的波长（nm）。
	LambdaNm float64
	// V 该波长下的归一化频率（由调用方传入）。
	V float64
	// DMat 材料色散，来自 Sellmeier 石英系数求导。
	DMat float64
	// DWg 波导色散，来自阶跃近似公式。
	DWg float64
	// DTotal 总色散 = DMat + DWg。
	DTotal float64
}

// DispersionTerms 返回 [DMat, DWg, DTotal] 三个分量，便于报告按列打印。
func (r Result) DispersionTerms() [3]float64 {
	return [3]float64{r.DMat, r.DWg, r.DTotal}
}

// Terms 返回色散项的名称，与 DispersionTerms 顺序一致。
func Terms() [3]string {
	return [3]string{"material", "waveguide", "total"}
}
