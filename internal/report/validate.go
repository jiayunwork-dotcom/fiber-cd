package report

import (
	"fmt"
	"io"
	"math"

	"fiber-cd/internal/model"
)

func PrintValidationReport(w io.Writer, cfg model.Config) error {
	fmt.Fprintf(w, "%s: valid step-index fiber\n", cfg.Description())
	fmt.Fprintf(w, "  n1 = %s\n", FmtFloat(cfg.N1, 6))
	fmt.Fprintf(w, "  n2 = %s\n", FmtFloat(cfg.N2, 6))
	fmt.Fprintf(w, "  core diameter = %s um\n", FmtFloat(cfg.CoreDiameterUm, 4))
	fmt.Fprintf(w, "  wavelength     = %s\n", FmtNm(cfg.WavelengthNm, 3))
	fmt.Fprintf(w, "  NA = %s\n", FmtFloat(modelNA(cfg), 6))
	return nil
}

func modelNA(cfg model.Config) float64 {
	return math.Sqrt(cfg.N1*cfg.N1 - cfg.N2*cfg.N2)
}
