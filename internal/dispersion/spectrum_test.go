package dispersion

import (
	"math"
	"testing"
)

func TestSpectrumAt(t *testing.T) {
	cfg := smfCfg()
	sp, err := SpectrumAt(cfg, 1310.0)
	if err != nil {
		t.Fatalf("SpectrumAt: %v", err)
	}
	if math.Abs(sp.Index-1.446804) > 1e-5 {
		t.Errorf("n(1310) = %v, want ≈ 1.446804", sp.Index)
	}
	if math.Abs(sp.DMat-3.581818) > 1e-3 {
		t.Errorf("D_mat = %v, want ≈ 3.581818", sp.DMat)
	}
	if math.Abs(sp.DTotal-(sp.DMat+sp.DWg)) > 1e-9 {
		t.Errorf("D_tot = %v, want D_mat + D_wg", sp.DTotal)
	}
}

func TestSpectrumAtRejectsBadWavelength(t *testing.T) {
	if _, err := SpectrumAt(smfCfg(), 0); err == nil {
		t.Error("expected error for zero wavelength, got nil")
	}
}
