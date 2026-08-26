package report

import (
	"bytes"
	"strings"
	"testing"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/waveguide"
)

func example() (waveguide.Result, dispersion.Result) {
	cfg := model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		panic(err)
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		panic(err)
	}
	return wr, dr
}

func TestModeReportContainsNumbers(t *testing.T) {
	wr, dr := example()
	var buf bytes.Buffer
	if err := PrintModeReport(&buf, wr, dr); err != nil {
		t.Fatalf("PrintModeReport: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"NA = 0.104076",
		"V = 2.24631",
		"single-mode",
		"1223.56 nm",
		"D_tot = -0.123",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q:\n%s", want, out)
		}
	}
}

func TestSweepTableMarksBoundary(t *testing.T) {
	cfg := model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
	sr, err := waveguide.Sweep(cfg, 1190, 1260, 8)
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	var buf bytes.Buffer
	if err := PrintSweepTable(&buf, cfg.Description(), sr); err != nil {
		t.Fatalf("PrintSweepTable: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "multi-mode") {
		t.Errorf("sweep table must show multi-mode rows below cutoff:\n%s", out)
	}
	if !strings.Contains(out, "single-mode") {
		t.Errorf("sweep table must show single-mode rows above cutoff:\n%s", out)
	}
	if !strings.Contains(out, "cutoff wavelength = 1223.56 nm") {
		t.Errorf("sweep table must report cutoff wavelength:\n%s", out)
	}
}

func TestValidationReportEchoesParams(t *testing.T) {
	cfg := model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
	var buf bytes.Buffer
	if err := PrintValidationReport(&buf, cfg); err != nil {
		t.Fatalf("PrintValidationReport: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"valid step-index fiber", "n1 = 1.465600", "wavelength     = 1310.000 nm"} {
		if !strings.Contains(out, want) {
			t.Errorf("validation report missing %q:\n%s", want, out)
		}
	}
}

func TestFormatHelpers(t *testing.T) {
	if got := FmtSigned(-3.704398, 3); got != "-3.704" {
		t.Errorf("FmtSigned(-3.704398,3) = %q, want %q", got, "-3.704")
	}
	if got := FmtPercent(0.00252456, 4); got != "0.2525 %" {
		t.Errorf("FmtPercent = %q, want %q", got, "0.2525 %")
	}
	if got := Truncate("abcdef", 4); got != "abc…" {
		t.Errorf("Truncate = %q, want %q", got, "abc…")
	}
}
