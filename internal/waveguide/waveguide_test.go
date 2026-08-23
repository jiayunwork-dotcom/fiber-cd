package waveguide

import (
	"math"
	"testing"

	"fiber-cd/internal/model"
)

func smfCfg() model.Config {
	return model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
}

const (
	smfNA      = 0.1040756936
	smfV       = 2.2463136718
	smfCutoff  = 1223.563788 // nm
	smfDMat    = 3.581818
	smfDWg     = -3.704398
	smfDTotal  = -0.122580
	smfZeroDNm = 1311.368
)

func TestNAFormulaMatchesDefinition(t *testing.T) {
	got := NumericalAperture(1.4656, 1.4619)
	if math.Abs(got-smfNA) > 1e-8 {
		t.Errorf("NA = %v, want ≈ %v", got, smfNA)
	}
	if got := NumericalAperture(1.5, 1.5); got != 0 {
		t.Errorf("NA for n1==n2 should be 0, got %v", got)
	}
}

func TestVNumberUsesRadiusNotDiameter(t *testing.T) {
	cfg := smfCfg()
	na := NumericalAperture(cfg.N1, cfg.N2)

	vRadius := VNumber(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	if math.Abs(vRadius-smfV) > 1e-6 {
		t.Errorf("V with radius a = %v, want ≈ %v (using diameter would give %v)", vRadius, smfV, 2*smfV)
	}

	vDiameter := VNumber(cfg.CoreDiameterUm*1e-6, na, cfg.WavelengthM())
	if math.Abs(vDiameter-2*smfV) > 1e-6 {
		t.Errorf("diameter-as-radius V = %v, want ≈ %v (2×)", vDiameter, 2*smfV)
	}

	cutoff := CutoffWavelengthNm(cfg.CoreRadiusM(), na)
	if math.Abs(cutoff-smfCutoff) > 1e-2 {
		t.Errorf("cutoff = %v nm, want ≈ %v nm (diameter-as-radius would give %v nm)", cutoff, smfCutoff, 2*smfCutoff)
	}

	res, err := Analyze(cfg)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if res.Mode != SingleMode {
		t.Errorf("smf-1310 must be single-mode, got %v", res.Mode)
	}
}

func TestCoreDiameterDoublingDoublesV(t *testing.T) {
	base := smfCfg()
	doubled := base.Clone()
	doubled.CoreDiameterUm = base.CoreDiameterUm * 2

	na := NumericalAperture(base.N1, base.N2)
	v1 := VNumber(base.CoreRadiusM(), na, base.WavelengthM())
	v2 := VNumber(doubled.CoreRadiusM(), na, doubled.WavelengthM())

	if ratio := v2 / v1; math.Abs(ratio-2.0) > 1e-9 {
		t.Errorf("doubling core diameter should double V, ratio = %v", ratio)
	}

	r1, _ := Analyze(base)
	r2, _ := Analyze(doubled)
	if r1.Mode != SingleMode {
		t.Errorf("base fiber should be single-mode, got %v", r1.Mode)
	}
	if r2.Mode != MultiMode {
		t.Errorf("doubled-diameter fiber at same wavelength should be multi-mode, got %v", r2.Mode)
	}
}

func TestWavelengthDoublingHalvesV(t *testing.T) {
	base := smfCfg()
	long := base.Clone()
	long.WavelengthNm = base.WavelengthNm * 2

	na := NumericalAperture(base.N1, base.N2)
	v1 := VNumber(base.CoreRadiusM(), na, base.WavelengthM())
	v2 := VNumber(long.CoreRadiusM(), na, long.WavelengthM())

	if ratio := v2 / v1; math.Abs(ratio-0.5) > 1e-9 {
		t.Errorf("doubling wavelength should halve V, ratio = %v", ratio)
	}
}

func TestReduceDeltaLowersNAAndV(t *testing.T) {
	base := smfCfg()
	tighter := base.Clone()
	tighter.N2 = 1.4640 // n1−n2 从 0.0037 缩到 0.0016

	naBase := NumericalAperture(base.N1, base.N2)
	naTight := NumericalAperture(tighter.N1, tighter.N2)
	if naTight >= naBase {
		t.Errorf("reduced delta must lower NA: %v vs %v", naTight, naBase)
	}

	vBase := VNumber(base.CoreRadiusM(), naBase, base.WavelengthM())
	vTight := VNumber(tighter.CoreRadiusM(), naTight, tighter.WavelengthM())
	if vTight >= vBase {
		t.Errorf("reduced delta must lower V: %v vs %v", vTight, vBase)
	}

	if math.Abs(naTight-0.06846430) > 1e-6 {
		t.Errorf("reduced-delta NA = %v, want ≈ 0.06846430", naTight)
	}
	if math.Abs(vTight-1.477696) > 1e-4 {
		t.Errorf("reduced-delta V = %v, want ≈ 1.477696", vTight)
	}
}

func TestSingleModeBoundaryAtCutoff(t *testing.T) {
	cases := []struct {
		v    float64
		want ModeStatus
	}{
		{v: 1.0, want: SingleMode},
		{v: 2.4049, want: SingleMode},
		{v: 2.405, want: MultiMode},
		{v: 2.4051, want: MultiMode},
		{v: 3.5, want: MultiMode},
	}
	for _, tc := range cases {
		if got := Classify(tc.v); got != tc.want {
			t.Errorf("Classify(%v) = %v, want %v", tc.v, got, tc.want)
		}
	}
}

func TestAnalyzeReportsExampleSingleMode(t *testing.T) {
	res, err := Analyze(smfCfg())
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if res.Mode != SingleMode {
		t.Errorf("smf-1310 at 1310 nm must be single-mode, got %v", res.Mode)
	}
	if math.Abs(res.V-smfV) > 1e-4 {
		t.Errorf("V = %v, want ≈ %v", res.V, smfV)
	}
	if math.Abs(res.CutoffWavelengthNm-smfCutoff) > 1.0 {
		t.Errorf("cutoff = %v nm, want ≈ %v nm", res.CutoffWavelengthNm, smfCutoff)
	}
}

func TestAnalyzeRejectsInvalidConfig(t *testing.T) {
	cfg := model.Config{N1: 1.4, N2: 1.5, CoreDiameterUm: 9, WavelengthNm: 1310}
	if _, err := Analyze(cfg); err == nil {
		t.Error("expected error for denser cladding, got nil")
	}
}

func TestSweepFlagsMultimodeBelowCutoff(t *testing.T) {
	sr, err := Sweep(smfCfg(), 1190, 1260, 8)
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	if len(sr.Steps) != 8 {
		t.Fatalf("steps = %d, want 8", len(sr.Steps))
	}
	for i := 0; i < 4; i++ {
		if sr.Steps[i].Mode != MultiMode {
			t.Errorf("step %d (%.0f nm, below cutoff) should be multi-mode, got %v",
				i, sr.Steps[i].WavelengthNm, sr.Steps[i].Mode)
		}
	}
	for i := 4; i < len(sr.Steps); i++ {
		if sr.Steps[i].Mode != SingleMode {
			t.Errorf("step %d (%.0f nm, above cutoff) should be single-mode, got %v",
				i, sr.Steps[i].WavelengthNm, sr.Steps[i].Mode)
		}
	}
	if !sr.CutoffFound {
		t.Error("cutoff should be inside the sweep range")
	}
	if math.Abs(sr.CutoffNm-smfCutoff) > 1e-2 {
		t.Errorf("sweep cutoff = %v nm, want ≈ %v nm", sr.CutoffNm, smfCutoff)
	}
}

func TestSweepZeroDispersion(t *testing.T) {
	sr, err := Sweep(smfCfg(), 1260, 1360, 21)
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	if !sr.ZeroDFound {
		t.Fatal("zero dispersion should be found in 1260–1360 nm")
	}
	if math.Abs(sr.ZeroDNm-smfZeroDNm) > 2.0 {
		t.Errorf("zero-D = %v nm, want ≈ %v nm", sr.ZeroDNm, smfZeroDNm)
	}
}

func TestSweepRejectsBadRange(t *testing.T) {
	if _, err := Sweep(smfCfg(), 0, 100, 10); err == nil {
		t.Error("expected error for nonpositive start, got nil")
	}
	if _, err := Sweep(smfCfg(), 1400, 1300, 10); err == nil {
		t.Error("expected error for stop < start, got nil")
	}
	if _, err := Sweep(smfCfg(), 1200, 1300, 1); err == nil {
		t.Error("expected error for steps < 2, got nil")
	}
}
