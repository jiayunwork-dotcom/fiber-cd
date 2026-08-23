package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/waveguide"
)

// PrintSpecReport 打印 spec 子命令的完整表征报告：mode 的导波量与
// 色散量之上，追加光斑尺寸、有效面积、模式数、群折射率与色散斜率。
func PrintSpecReport(w io.Writer, wr waveguide.Result, dr dispersion.Result, sr dispersion.SlopeResult) error {
	cfg := wr.Config

	fmt.Fprintf(w, "%s\n", cfg.Description())
	fmt.Fprintf(w, "  n1 = %s   n2 = %s   delta = %s\n",
		FmtFloat(cfg.N1, 6), FmtFloat(cfg.N2, 6), FmtPercent(wr.Delta, 4))
	fmt.Fprintf(w, "  core diameter = %s um   (radius a = %s um)\n",
		FmtFloat(cfg.CoreDiameterUm, 4), FmtFloat(cfg.RadiusUm(), 4))
	fmt.Fprintf(w, "  wavelength     = %s\n\n", FmtNm(cfg.WavelengthNm, 3))

	fmt.Fprintf(w, "Guided-wave quantities:\n")
	fmt.Fprintf(w, "  NA       = %s\n", FmtFloat(wr.NA, 6))
	fmt.Fprintf(w, "  V        = %s\n", FmtFloat(wr.V, 5))
	fmt.Fprintf(w, "  cutoff   = %s\n", FmtNm(wr.CutoffWavelengthNm, 2))
	fmt.Fprintf(w, "  mode     = %s\n", ModeStatusLabel(wr.Mode))
	fmt.Fprintf(w, "  modes    = %s\n", modeCountText(wr))
	fmt.Fprintf(w, "  MFD      = %s um\n", FmtFloat(waveguide.ModeFieldDiameterUm(cfg.RadiusUm(), wr.V), 3))
	fmt.Fprintf(w, "  Aeff     = %s um^2\n\n", FmtFloat(waveguide.EffectiveAreaUm2(cfg.RadiusUm(), wr.V), 3))

	fmt.Fprintf(w, "Material spectral data (fused silica):\n")
	fmt.Fprintf(w, "  n(λ)      = %s\n", FmtFloat(dispersion.Silica().IndexUm(cfg.WavelengthUm()), 6))
	fmt.Fprintf(w, "  N_g(λ)    = %s\n", FmtFloat(dispersion.GroupIndex(cfg.WavelengthUm()), 6))
	fmt.Fprintf(w, "  τ         = %s us/km\n\n", FmtFloat(dispersion.GroupDelayPerKm(cfg.WavelengthUm()), 3))

	fmt.Fprintf(w, "Dispersion (ps/(nm*km)) at %s:\n", FmtNm(dr.LambdaNm, 3))
	fmt.Fprintf(w, "  D_mat = %s\n", FmtSigned(dr.DMat, 3))
	fmt.Fprintf(w, "  D_wg  = %s\n", FmtSigned(dr.DWg, 3))
	fmt.Fprintf(w, "  D_tot = %s\n", FmtSigned(dr.DTotal, 3))
	fmt.Fprintf(w, "  S_tot = %s ps/(nm^2*km)\n", FmtSigned(sr.STotal, 4))
	fmt.Fprintf(w, "  S_mat = %s ps/(nm^2*km)\n", FmtSigned(sr.SMaterial, 4))
	fmt.Fprintf(w, "  S_wg  = %s ps/(nm^2*km)\n", FmtSigned(sr.SWaveguide, 4))
	return nil
}

// modeCountText 复用 waveguide 的模式数描述，避免 report 重复实现。
func modeCountText(wr waveguide.Result) string {
	return waveguide.ModeCountDescription(wr.V)
}
