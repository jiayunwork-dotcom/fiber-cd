package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/report"
	"fiber-cd/internal/server"
	"fiber-cd/internal/waveguide"
)

const usageText = `usage: fiber-cd <subcommand> <config.json> [args]

Subcommands:
  mode <config.json>                                  single-shot: NA, V, mode, D_tot
  spec <config.json>                                  full characterization (mode + dispersion + slope)
  design <config.json>                                single-mode design limits (max diameter / NA)
  probe <config.json> <wavelength_nm>                 status at an arbitrary wavelength
  boundary <config.json>                              dispersion exactly at the V=2.405 cutoff
  sweep <config.json> [start_nm stop_nm [steps]]      wavelength sweep table
  dump <config.json>                                  re-emit the normalized config as JSON
  validate <config.json>                              validate a fiber config only
  version                                             print version
  help                                                show this usage

Anywhere a <config.json> path is expected, "-" reads the config from stdin.

Examples:
  fiber-cd mode example/smf-1310.json
  fiber-cd spec example/smf-1310.json
  fiber-cd design example/smf-1310.json
  fiber-cd probe example/smf-1310.json 1550.0
  fiber-cd boundary example/smf-1310.json
  fiber-cd sweep example/smf-1310.json 1260 1620 41
  fiber-cd dump example/smf-1310.json
  fiber-cd validate example/smf-1310.json
`

const version = "fiber-cd 1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, usageText)
		return 2
	}
	switch args[0] {
	case "mode":
		return runMode(args[1:], stdout, stderr)
	case "-http", "--http":
		return runHTTP(args[1:], stdout, stderr)
	case "spec":
		return runSpec(args[1:], stdout, stderr)
	case "design":
		return runDesign(args[1:], stdout, stderr)
	case "probe":
		return runProbe(args[1:], stdout, stderr)
	case "boundary":
		return runBoundary(args[1:], stdout, stderr)
	case "dump":
		return runDump(args[1:], stdout, stderr)
	case "sweep":
		return runSweep(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "version", "-v":
		fmt.Fprintln(stdout, version)
		return 0
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usageText)
		return 0
	default:
		fmt.Fprintf(stderr, "fiber-cd: unknown subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, usageText)
		return 2
	}
}

func loadConfigArg(path string) (model.Config, error) {
	if path == "-" {
		return model.Decode(os.Stdin)
	}
	return model.LoadFile(path)
}

func runSpec(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd spec: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	sl, err := dispersion.SlopeAtOperating(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintSpecReport(stdout, wr, dr, sl); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runDesign(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd design: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintDesignReport(stdout, wr); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runProbe(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintf(stderr, "fiber-cd probe: usage: probe <config.json> <wavelength_nm>\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	wavelengthNm, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd probe: invalid wavelength %q: %v\n", args[1], err)
		return 1
	}
	pr, err := waveguide.Probe(cfg, wavelengthNm)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintProbeReport(stdout, pr); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runBoundary(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd boundary: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	cutoffM := wr.CutoffWavelengthNm * 1e-9
	dr, err := dispersion.Compose(cfg, waveguide.CutoffV, cutoffM)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintBoundaryReport(stdout, cfg, wr, dr); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runDump(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd dump: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := model.Validate(cfg); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	data, err := cfg.MarshalIndent()
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd dump: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, string(data))
	return 0
}

func runMode(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd mode: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintModeReport(stdout, wr, dr); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runSweep(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 || len(args) > 4 {
		fmt.Fprintf(stderr, "fiber-cd sweep: usage: sweep <config.json> [start_nm stop_nm [steps]]\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}

	start, stop, steps := cfg.SweepRange()
	if len(args) >= 3 {
		start, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(stderr, "fiber-cd sweep: invalid start wavelength %q: %v\n", args[1], err)
			return 1
		}
		stop, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			fmt.Fprintf(stderr, "fiber-cd sweep: invalid stop wavelength %q: %v\n", args[2], err)
			return 1
		}
	}
	if len(args) == 4 {
		steps, err = strconv.Atoi(args[3])
		if err != nil {
			fmt.Fprintf(stderr, "fiber-cd sweep: invalid step count %q: %v\n", args[3], err)
			return 1
		}
	}

	sr, err := waveguide.Sweep(cfg, start, stop, steps)
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintSweepTable(stdout, cfg.Description(), sr); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "fiber-cd validate: expected exactly one config path\n\n")
		fmt.Fprint(stderr, usageText)
		return 2
	}
	cfg, err := loadConfigArg(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := model.Validate(cfg); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	if err := report.PrintValidationReport(stdout, cfg); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: %v\n", err)
		return 1
	}
	return 0
}

func runHTTP(args []string, stdout, stderr io.Writer) int {
	addr := ":8080"
	if len(args) > 0 {
		addr = args[0]
	}
	if len(args) > 1 {
		fmt.Fprintf(stderr, "fiber-cd: -http accepts at most one listen address\n")
		return 2
	}
	fmt.Fprintf(stdout, "fiber-cd HTTP on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, server.New("web", "example/smf-1310.json")); err != nil {
		fmt.Fprintf(stderr, "fiber-cd: HTTP server: %v\n", err)
		return 1
	}
	return 0
}
