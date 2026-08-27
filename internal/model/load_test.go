package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfig(t *testing.T) {
	data := []byte(`{
	  "name": "probe",
	  "n1": 1.4656,
	  "n2": 1.4619,
	  "core_diameter_um": 9.0,
	  "wavelength_nm": 1310.0
	}`)
	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if cfg.N1 != 1.4656 || cfg.N2 != 1.4619 {
		t.Errorf("indices parsed wrong: n1=%v n2=%v", cfg.N1, cfg.N2)
	}
	if cfg.CoreDiameterUm != 9.0 || cfg.WavelengthNm != 1310.0 {
		t.Errorf("geometry parsed wrong: d=%v lam=%v", cfg.CoreDiameterUm, cfg.WavelengthNm)
	}
	if cfg.Name != "probe" {
		t.Errorf("name parsed wrong: %q", cfg.Name)
	}
}

func TestParseConfigRejectsBadJSON(t *testing.T) {
	if _, err := Parse([]byte(`{"n1": }`)); err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

func TestLoadFileMissing(t *testing.T) {
	if _, err := LoadFile(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	data := `{"n1":1.48,"n2":1.465,"core_diameter_um":62.5,"wavelength_nm":850}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}
	if cfg.CoreDiameterUm != 62.5 || cfg.WavelengthNm != 850 {
		t.Errorf("round trip mismatch: %+v", cfg)
	}
}

func TestSweepRangeDefaults(t *testing.T) {
	cfg := Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9, WavelengthNm: 1310}
	start, stop, steps := cfg.SweepRange()
	if start != DefaultSweepStartNm || stop != DefaultSweepStopNm || steps != DefaultSweepSteps {
		t.Errorf("defaults mismatch: start=%v stop=%v steps=%d", start, stop, steps)
	}
}

func TestSweepRangeFromConfig(t *testing.T) {
	cfg := Config{
		N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9, WavelengthNm: 1310,
		SweepStartNm: 1260, SweepStopNm: 1620, SweepSteps: 41,
	}
	start, stop, steps := cfg.SweepRange()
	if start != 1260 || stop != 1620 || steps != 41 {
		t.Errorf("config sweep range ignored: start=%v stop=%v steps=%d", start, stop, steps)
	}
}
