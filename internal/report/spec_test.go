package report

import (
	"bytes"
	"strings"
	"testing"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/waveguide"
)

func mustCfg() model.Config {
	return model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
}

func TestSpecReportContainsSlope(t *testing.T) {
	cfg := mustCfg()
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		t.Fatal(err)
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		t.Fatal(err)
	}
	sl, err := dispersion.SlopeAtOperating(cfg)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := PrintSpecReport(&buf, wr, dr, sl); err != nil {
		t.Fatalf("PrintSpecReport: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"S_tot = +0.0897", "N_g(λ)", "MFD", "Aeff"} {
		if !strings.Contains(out, want) {
			t.Errorf("spec report missing %q:\n%s", want, out)
		}
	}
}

func TestDesignReportContainsLimits(t *testing.T) {
	wr, err := waveguide.Analyze(mustCfg())
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := PrintDesignReport(&buf, wr); err != nil {
		t.Fatalf("PrintDesignReport: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"max core diameter = 9.636 um", "max NA", "verdict: single-mode"} {
		if !strings.Contains(out, want) {
			t.Errorf("design report missing %q:\n%s", want, out)
		}
	}
}
