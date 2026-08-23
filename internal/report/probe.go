package report

import (
	"fmt"
	"io"

	"fiber-cd/internal/waveguide"
)

// PrintProbeReport 打印 probe 子命令的结果：光纤在指定波长下的
// V、模式与三项色散。
func PrintProbeReport(w io.Writer, pr waveguide.ProbeResult) error {
	fmt.Fprintf(w, "%s at %s\n", pr.Config.Description(), FmtNm(pr.WavelengthNm, 3))
	fmt.Fprintf(w, "  V      = %s\n", FmtFloat(pr.V, 5))
	fmt.Fprintf(w, "  cutoff = %s\n", FmtNm(pr.CutoffNm, 2))
	fmt.Fprintf(w, "  mode   = %s\n", ModeStatusLabel(pr.Mode))
	fmt.Fprintf(w, "  D_mat  = %s ps/(nm*km)\n", FmtSigned(pr.Disp.DMat, 3))
	fmt.Fprintf(w, "  D_wg   = %s ps/(nm*km)\n", FmtSigned(pr.Disp.DWg, 3))
	fmt.Fprintf(w, "  D_tot  = %s ps/(nm*km)\n", FmtSigned(pr.Disp.DTotal, 3))
	return nil
}
