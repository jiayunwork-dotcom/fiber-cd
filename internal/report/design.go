package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/waveguide"
)

// PrintDesignReport 打印 design 子命令的逆向设计结论：当前配置的
// 单模裕量、保持单模允许的最大芯径与最大 NA。
func PrintDesignReport(w io.Writer, wr waveguide.Result) error {
	cfg := wr.Config

	fmt.Fprintf(w, "%s: single-mode design check\n", cfg.Description())
	fmt.Fprintf(w, "  wavelength     = %s\n", FmtNm(cfg.WavelengthNm, 3))
	fmt.Fprintf(w, "  NA             = %s\n", FmtFloat(wr.NA, 6))
	fmt.Fprintf(w, "  core diameter  = %s um\n", FmtFloat(cfg.CoreDiameterUm, 4))
	fmt.Fprintf(w, "  V              = %s\n", FmtFloat(wr.V, 5))
	fmt.Fprintf(w, "  margin         = %s  (2.405 − V)\n", FmtSigned(waveguide.SingleModeMargin(wr.V), 5))
	fmt.Fprintf(w, "  margin %%       = %s %%\n", FmtFloat(waveguide.SingleModeMarginPercent(wr.V), 2))

	maxD := waveguide.MaxSingleModeDiameterUm(wr.NA, cfg.WavelengthNm)
	maxNA := waveguide.MaxSingleModeNA(cfg.CoreDiameterUm, cfg.WavelengthNm)
	fmt.Fprintf(w, "\nKeeping V ≤ 2.405 at %s:\n", FmtNm(cfg.WavelengthNm, 3))
	fmt.Fprintf(w, "  max core diameter = %s um\n", FmtFloat(maxD, 3))
	fmt.Fprintf(w, "  max NA            = %s\n", FmtFloat(maxNA, 6))

	if wr.Mode == waveguide.SingleMode {
		fmt.Fprintf(w, "\nverdict: single-mode (within %.1f %% of the V=2.405 limit)\n",
			waveguide.SingleModeMarginPercent(wr.V))
		fmt.Fprintf(w, "headroom: core diameter may grow by up to %.1f %% before multi-mode\n",
			(maxD-cfg.CoreDiameterUm)/cfg.CoreDiameterUm*100)
	} else {
		// 已经多模：给出把芯径缩回多少能回到单模。
		targetD := waveguide.MaxSingleModeDiameterUm(wr.NA, cfg.WavelengthNm)
		fmt.Fprintf(w, "\nverdict: multi-mode at this wavelength\n")
		fmt.Fprintf(w, "to return to single-mode: core diameter must shrink to %.3f um\n", targetD)
	}
	return nil
}
