package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const examplePath = "example/smf-1310.json"

func TestCLIModeCommand(t *testing.T) {
	if _, err := os.Stat(examplePath); err != nil {
		t.Skipf("example not present in %s: %v", examplePath, err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"mode", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"NA = 0.104076", "single-mode", "D_tot"} {
		if !strings.Contains(out, want) {
			t.Errorf("mode output missing %q:\n%s", want, out)
		}
	}
	if stderr.Len() != 0 {
		t.Errorf("mode output leaked to stderr: %s", stderr.String())
	}
}

func TestCLISweepCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"sweep", examplePath, "1260", "1360", "11"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "zero-dispersion wavelength") {
		t.Errorf("sweep output missing zero-dispersion conclusion:\n%s", stdout.String())
	}
}

func TestCLIRejectsInvalidInput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	bad := `{"n1":1.4,"n2":1.5,"core_diameter_um":9,"wavelength_nm":1310}`
	if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"mode", path}, &stdout, &stderr)
	if code == 0 {
		t.Error("exit code = 0, want non-zero for denser cladding")
	}
	if !strings.Contains(stderr.String(), "cladding") {
		t.Errorf("stderr should explain the invalid input, got: %s", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("invalid input must not print a report, got: %s", stdout.String())
	}
}

func TestCLIRejectsMissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"mode", filepath.Join(t.TempDir(), "nope.json")}, &stdout, &stderr)
	if code == 0 {
		t.Error("exit code = 0, want non-zero for missing file")
	}
	if stderr.Len() == 0 {
		t.Error("expected an error message on stderr")
	}
}

func TestCLIUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	if code != 2 {
		t.Errorf("empty args exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "usage:") {
		t.Errorf("stderr should contain usage, got: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"frobnicate", "x"}, &stdout, &stderr)
	if code != 2 {
		t.Errorf("unknown subcommand exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown subcommand") {
		t.Errorf("stderr should name the unknown subcommand, got: %s", stderr.String())
	}
}

func TestCLIValidateCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"validate", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "NA = 0.104076") {
		t.Errorf("validate output missing NA:\n%s", stdout.String())
	}
}

func TestCLIHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"help"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("help exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "usage:") {
		t.Errorf("help should print usage to stdout, got: %s", stdout.String())
	}
}

func TestCLISpecCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"spec", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"S_tot", "MFD", "N_g(λ)"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("spec output missing %q:\n%s", want, stdout.String())
		}
	}
}

func TestCLIDesignCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"design", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"max core diameter", "max NA", "verdict"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("design output missing %q:\n%s", want, stdout.String())
		}
	}
}

func TestCLIProbeCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"probe", examplePath, "1550.0"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"at 1550.000 nm", "single-mode", "D_tot"} {
		if !strings.Contains(out, want) {
			t.Errorf("probe output missing %q:\n%s", want, out)
		}
	}
}

func TestCLIBoundaryCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"boundary", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"2.4048", "1223.56 nm", "V = 2.405"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("boundary output missing %q:\n%s", want, stdout.String())
		}
	}
}

func TestCLIDumpCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"dump", examplePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"core_diameter_um": 9`) {
		t.Errorf("dump output missing core diameter:\n%s", stdout.String())
	}
}

func TestCLIStdinReadsConfig(t *testing.T) {
	var stdout, stderr bytes.Buffer
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdin
	os.Stdin = pr
	defer func() { os.Stdin = orig }()
	go func() {
		pw.Write([]byte(`{"n1":1.4656,"n2":1.4619,"core_diameter_um":9,"wavelength_nm":1310}`))
		pw.Close()
	}()
	code := run([]string{"mode", "-"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "single-mode") {
		t.Errorf("stdin mode output should be single-mode:\n%s", stdout.String())
	}
}

func TestCLIVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("version exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "fiber-cd") {
		t.Errorf("version output missing tool name: %s", stdout.String())
	}
}
