package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/waveguide"
)

func PrintModeReport(w io.Writer, wr waveguide.Result, dr dispersion.Result) error {
	cfg := wr.Config

	fmt.Fprintf(w, "%s\n", cfg.Description())
	fmt.Fprintf(w, "  n1 = %s   n2 = %s   delta = %s\n",
		FmtFloat(cfg.N1, 6), FmtFloat(cfg.N2, 6), FmtPercent(wr.Delta, 4))
	fmt.Fprintf(w, "  core diameter = %s um   (radius a = %s um)\n",
		FmtFloat(cfg.CoreDiameterUm, 4), FmtFloat(cfg.RadiusUm(), 4))
	fmt.Fprintf(w, "  wavelength     = %s\n\n", FmtNm(cfg.WavelengthNm, 3))

	fmt.Fprintf(w, "Numerical aperture  NA = %s\n", FmtFloat(wr.NA, 6))
	fmt.Fprintf(w, "Normalized frequency V = %s     (single-mode iff V < 2.405)\n", FmtFloat(wr.V, 5))
	fmt.Fprintf(w, "Cutoff wavelength        = %s\n", FmtNm(wr.CutoffWavelengthNm, 2))
	fmt.Fprintf(w, "Mode status              = %s\n\n", ModeStatusLabel(wr.Mode))

	fmt.Fprintf(w, "Chromatic dispersion (ps/(nm*km)) at %s:\n", FmtNm(dr.LambdaNm, 3))
	fmt.Fprintf(w, "  material  D_mat = %s\n", FmtSigned(dr.DMat, 3))
	fmt.Fprintf(w, "  waveguide D_wg  = %s\n", FmtSigned(dr.DWg, 3))
	fmt.Fprintf(w, "  total     D_tot = %s\n", FmtSigned(dr.DTotal, 3))
	fmt.Fprintf(w, "\nsummary: %s, NA=%s, V=%s, %s, D_tot=%s ps/(nm*km)\n",
		cfg.Description(), FmtFloat(wr.NA, 6), FmtFloat(wr.V, 5),
		ModeStatusLabel(wr.Mode), FmtSigned(dr.DTotal, 3))
	return nil
}
