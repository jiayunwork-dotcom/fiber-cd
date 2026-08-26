package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/waveguide"
)

func PrintBoundaryReport(w io.Writer, cfg model.Config, wr waveguide.Result, dr dispersion.Result) error {
	fmt.Fprintf(w, "%s: single-mode boundary (V = 2.405)\n", cfg.Description())
	fmt.Fprintf(w, "  cutoff wavelength = %s\n", FmtNm(wr.CutoffWavelengthNm, 2))
	fmt.Fprintf(w, "  cutoff V (J0 first root) = %s\n", FmtFloat(waveguide.CutoffVFromJ0(), 4))
	fmt.Fprintf(w, "  NA = %s\n\n", FmtFloat(wr.NA, 6))

	fmt.Fprintf(w, "  dispersion at cutoff:\n")
	fmt.Fprintf(w, "    D_mat = %s ps/(nm*km)\n", FmtSigned(dr.DMat, 3))
	fmt.Fprintf(w, "    D_wg  = %s ps/(nm*km)\n", FmtSigned(dr.DWg, 3))
	fmt.Fprintf(w, "    D_tot = %s ps/(nm*km)\n", FmtSigned(dr.DTotal, 3))
	return nil
}
