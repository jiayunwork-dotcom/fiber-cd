package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/waveguide"
)

// PrintSweepTable 打印 sweep 子命令的波长扫描表：每档波长对应 V、
// 模式状态、模式数、三项色散与斜率，表尾附单模边界、零色散波长与
// 色散最平坦点结论。模式发生翻转的档与色散跨零的档会打标记。
func PrintSweepTable(w io.Writer, desc string, sr waveguide.SweepResult) error {
	fmt.Fprintf(w, "%s\nwavelength sweep %s .. %s, %d steps\n\n",
		desc, FmtNm(sr.StartNm, 0), FmtNm(sr.StopNm, 0), len(sr.Steps))

	// 找到模式翻转档与零色散档，用于行尾标记。
	flipIdx := -1
	for i := 1; i < len(sr.Steps); i++ {
		if sr.Steps[i-1].Mode != sr.Steps[i].Mode {
			flipIdx = i
			break
		}
	}
	zeroIdx := -1
	for i := 1; i < len(sr.Steps); i++ {
		if (sr.Steps[i-1].Disp.DTotal < 0) != (sr.Steps[i].Disp.DTotal < 0) {
			zeroIdx = i
			break
		}
	}

	tbl := NewTable("lambda(nm)", "V", "mode", "modes", "D_mat", "D_wg", "D_tot", "S_tot")
	for i, s := range sr.Steps {
		mark := ""
		switch i {
		case flipIdx:
			mark = "  <- single-mode boundary"
		case zeroIdx:
			mark = "  <- zero-D crossing"
		}
		tbl.AddRow(
			fmt.Sprintf("%.2f", s.WavelengthNm),
			fmt.Sprintf("%.4f", s.V),
			ModeStatusLabel(s.Mode),
			fmt.Sprintf("%.0f", waveguide.EffectiveModeCount(s.V)),
			fmt.Sprintf("%+.2f", s.Disp.DMat),
			fmt.Sprintf("%+.2f", s.Disp.DWg),
			fmt.Sprintf("%+.2f", s.Disp.DTotal),
			fmt.Sprintf("%+.3f%s", s.Slope, mark),
		)
	}
	if err := tbl.Render(w); err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	if sr.CutoffFound {
		fmt.Fprintf(w, "single-mode boundary: cutoff wavelength = %s\n", FmtNm(sr.CutoffNm, 2))
	} else {
		fmt.Fprintf(w, "single-mode boundary: outside sweep range\n")
	}
	if lo, hi, ok := sr.SingleModeWindow(); ok {
		fmt.Fprintf(w, "single-mode window   : %s .. %s\n", FmtNm(lo, 2), FmtNm(hi, 2))
	} else {
		fmt.Fprintf(w, "single-mode window   : none in sweep range\n")
	}
	if sr.ZeroDFound {
		fmt.Fprintf(w, "zero-dispersion wavelength = %s\n", FmtNm(sr.ZeroDNm, 1))
	} else {
		fmt.Fprintf(w, "zero-dispersion: no crossing in sweep range\n")
	}
	if lam, d, ok := sr.MinAbsDispersion(); ok {
		fmt.Fprintf(w, "minimum |D_tot|     : %+.3f ps/(nm*km) at %s\n", d, FmtNm(lam, 2))
	}
	return nil
}
