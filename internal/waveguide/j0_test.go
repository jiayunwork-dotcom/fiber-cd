package waveguide

import (
	"math"
	"testing"
)

func TestJ0FirstRoot(t *testing.T) {
	root := J0FirstRoot()
	if math.Abs(root-2.404826) > 1e-4 {
		t.Errorf("J0 first root = %v, want ≈ 2.404826", root)
	}
	if math.Abs(root-CutoffV) > 0.001 {
		t.Errorf("J0 root %v should match CutoffV %v", root, CutoffV)
	}
}

func TestJ0KnownValues(t *testing.T) {
	cases := []struct {
		x    float64
		want float64
	}{
		{x: 0, want: 1.0},
		{x: 2.4048255577, want: 0.0},
		{x: 3.8317059702, want: -0.40276}, // J1 首根处 J0 为负
	}
	for _, tc := range cases {
		got := J0(tc.x)
		if math.Abs(got-tc.want) > 1e-3 {
			t.Errorf("J0(%v) = %v, want ≈ %v", tc.x, got, tc.want)
		}
	}
}

func TestSweepSingleModeWindow(t *testing.T) {
	sr, err := Sweep(smfCfg(), 1190, 1260, 8)
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	lo, hi, ok := sr.SingleModeWindow()
	if !ok {
		t.Fatal("expected a single-mode window in range")
	}
	if math.Abs(lo-1230) > 1e-6 {
		t.Errorf("window start = %v nm, want 1230 nm", lo)
	}
	if math.Abs(hi-1260) > 1e-6 {
		t.Errorf("window end = %v nm, want 1260 nm", hi)
	}
}

func TestMinAbsDispersion(t *testing.T) {
	sr, err := Sweep(smfCfg(), 1260, 1360, 41)
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	lam, d, ok := sr.MinAbsDispersion()
	if !ok {
		t.Fatal("expected a min-|D| point")
	}
	if math.Abs(d) > 0.5 {
		t.Errorf("min |D_tot| = %v at %v nm, want close to 0", d, lam)
	}
	if lam < 1300 || lam > 1325 {
		t.Errorf("min-|D| wavelength = %v nm, want ≈ 1311 nm", lam)
	}
}
