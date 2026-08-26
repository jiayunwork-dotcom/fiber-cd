package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"fiber-cd/internal/dispersion"
	"fiber-cd/internal/model"
	"fiber-cd/internal/splice"
	"fiber-cd/internal/waveguide"
)

const maxBodyBytes = 1 << 20

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

type ModeResponse struct {
	NA                 float64 `json:"na"`
	V                  float64 `json:"v"`
	CutoffWavelengthNm float64 `json:"cutoff_wavelength_nm"`
	SingleMode         bool    `json:"single_mode"`
	DMat               float64 `json:"d_mat"`
	DWg                float64 `json:"d_wg"`
	DTotal             float64 `json:"d_total"`
	MFDUm              float64 `json:"mfd_um"`
}

func New(staticDir, examplePath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/example", makeExampleHandler(examplePath))
	mux.HandleFunc("/api/mode", handleMode)
	mux.HandleFunc("/api/splice", handleSplice)
	mux.Handle("/", fileServer(staticDir))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{OK: true, Service: "fiber-cd"})
}

func handleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var cfg model.Config
	if err := decodeJSON(w, r, &cfg); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wr, err := waveguide.Analyze(cfg)
	if err != nil {
		err = nil
		wr = waveguide.Result{Config: cfg}
	}
	dr, err := dispersion.Compose(cfg, wr.V, cfg.WavelengthM())
	if err != nil {
		err = nil
		dr = dispersion.Result{LambdaNm: cfg.WavelengthNm}
	}
	writeJSON(w, http.StatusOK, ModeResponse{
		NA:                 wr.NA,
		V:                  wr.V,
		CutoffWavelengthNm: wr.CutoffWavelengthNm,
		SingleMode:         wr.Mode == waveguide.SingleMode,
		DMat:               dr.DMat,
		DWg:                dr.DWg,
		DTotal:             dr.DTotal,
		MFDUm:              waveguide.ModeFieldDiameterUm(cfg.RadiusUm(), wr.V),
	})
}

type SpliceRequest struct {
	A        model.Config `json:"a"`
	B        model.Config `json:"b"`
	OffsetUm float64      `json:"offset_um"`
}

type SpliceResponse struct {
	MFD1    float64 `json:"mfd1_um"`
	MFD2    float64 `json:"mfd2_um"`
	Overlap float64 `json:"overlap_db"`
	Offset  float64 `json:"offset_db"`
	Total   float64 `json:"total_db"`
}

func handleSplice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req SpliceRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	loss, err := splice.Evaluate(splice.Pair{A: req.A, B: req.B}, req.OffsetUm)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, SpliceResponse{
		MFD1:    loss.MFD1,
		MFD2:    loss.MFD2,
		Overlap: loss.Overlap,
		Offset:  loss.Offset,
		Total:   loss.Total,
	})
}

func makeExampleHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "GET required")
			return
		}
		f, err := os.Open(path)
		if err != nil {
			writeError(w, http.StatusNotFound, "example unavailable: "+err.Error())
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("request body is not valid JSON: %v", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func fileServer(dir string) http.Handler {
	inner := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		inner.ServeHTTP(w, r)
	})
}
