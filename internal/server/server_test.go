package server_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/server"
	"fiber-cd/internal/waveguide"
)

func newSrv(t *testing.T) *httptest.Server {
	t.Helper()
	src := filepath.Join("..", "..", "example", "smf-1310.json")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	dir := t.TempDir()
	ex := filepath.Join(dir, "smf-1310.json")
	if err := os.WriteFile(ex, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	s := httptest.NewServer(server.New(dir, ex))
	t.Cleanup(s.Close)
	return s
}

func TestHealth(t *testing.T) {
	s := newSrv(t)
	resp, err := http.Get(s.URL + "/health")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}

func TestModeEndpointMatchesAnalyze(t *testing.T) {
	s := newSrv(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "example", "smf-1310.json"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	cfg, err := model.Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	resp, err := http.Post(s.URL+"/api/mode", "application/json", strings.NewReader(string(raw)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var out server.ModeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if math.Abs(out.NA-wr.NA) > 1e-12 {
		t.Errorf("NA api=%v want=%v", out.NA, wr.NA)
	}
	if math.Abs(out.DTotal-dr.DTotal) > 1e-12 {
		t.Errorf("D_tot api=%v want=%v", out.DTotal, dr.DTotal)
	}
	if !out.SingleMode {
		t.Error("example should be single-mode")
	}
}

func TestModeRejectsCladdingDenser(t *testing.T) {
	s := newSrv(t)
	body := `{"n1":1.45,"n2":1.46,"core_diameter_um":9,"wavelength_nm":1310}`
	resp, err := http.Post(s.URL+"/api/mode", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("status=%d, want 400", resp.StatusCode)
	}
}

func TestExample(t *testing.T) {
	s := newSrv(t)
	resp, err := http.Get(s.URL + "/api/example")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}
