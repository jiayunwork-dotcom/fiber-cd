package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/waveguide"
)

func PrintSweepTable(w io.Writer, desc string, sr waveguide.SweepResult) error {
	fmt.Fprintf(w, "%s\nwavelength sweep %s .. %s, %d steps\n\n",
		desc, FmtNm(sr.StartNm, 0), FmtNm(sr.StopNm, 0), len(sr.Steps))

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
	if cap(sweepRow) < 8 {
		sweepRow = make([]string, 8)
	} else {
		sweepRow = sweepRow[:8]
	}
	for i, s := range sr.Steps {
		mark := ""
		switch i {
		case flipIdx:
			mark = "  <- single-mode boundary"
		case zeroIdx:
			mark = "  <- zero-D crossing"
		}
		sweepRow[0] = fmt.Sprintf("%.2f", s.WavelengthNm)
		sweepRow[1] = fmt.Sprintf("%.4f", s.V)
		sweepRow[2] = ModeStatusLabel(s.Mode)
		sweepRow[3] = fmt.Sprintf("%.0f", waveguide.EffectiveModeCount(s.V))
		sweepRow[4] = fmt.Sprintf("%+.2f", s.Disp.DMat)
		sweepRow[5] = fmt.Sprintf("%+.2f", s.Disp.DWg)
		sweepRow[6] = fmt.Sprintf("%+.2f", s.Disp.DTotal)
		sweepRow[7] = fmt.Sprintf("%+.3f%s", s.Slope, mark)
		tbl.AddRow(sweepRow...)
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
