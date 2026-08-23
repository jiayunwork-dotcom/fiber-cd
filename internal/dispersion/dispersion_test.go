package dispersion

import (
	"errors"
	"math"
	"testing"

	"fiber-cd/internal/model"
)

func smfCfg() model.Config {
	return model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
}

func TestSellmeierIndexMatchesSilica(t *testing.T) {
	s := Silica()
	if got := s.IndexUm(1.31); math.Abs(got-1.446804) > 1e-5 {
		t.Errorf("n(1310nm) = %v, want ≈ 1.446804", got)
	}
	if got := s.IndexUm(1.55); math.Abs(got-1.444024) > 1e-5 {
		t.Errorf("n(1550nm) = %v, want ≈ 1.444024", got)
	}
}

func TestMaterialDispersionSignFlip(t *testing.T) {
	dm1250 := MaterialDispersionUm(1.25)
	if dm1250 >= 0 {
		t.Errorf("D_mat(1250nm) = %v, want negative", dm1250)
	}
	dm1600 := MaterialDispersionUm(1.60)
	if dm1600 <= 0 {
		t.Errorf("D_mat(1600nm) = %v, want positive", dm1600)
	}
	if math.Abs(dm1250-(-2.332279)) > 1e-3 {
		t.Errorf("D_mat(1250nm) = %v, want ≈ -2.332279", dm1250)
	}
	if math.Abs(dm1600-25.071598) > 1e-3 {
		t.Errorf("D_mat(1600nm) = %v, want ≈ 25.071598", dm1600)
	}
}

func TestMaterialDispersionZero(t *testing.T) {
	lo, hi := 1.20, 1.40
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if MaterialDispersionUm(mid) < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	zero := (lo + hi) / 2 * 1e3
	if zero < 1270 || zero > 1276 {
		t.Errorf("silica zero material dispersion = %v nm, want ≈ 1272.8 nm", zero)
	}
}

func TestTotalDispersionIsSum(t *testing.T) {
	cfg := smfCfg()
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	res, err := Compose(cfg, v, cfg.WavelengthM())
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	want := res.DMat + res.DWg
	if math.Abs(res.DTotal-want) > 1e-9 {
		t.Errorf("D_tot = %v, want D_mat+D_wg = %v", res.DTotal, want)
	}
}

func TestWaveguideDispersionSignAndValue(t *testing.T) {
	cfg := smfCfg()
	na := NumericalAperture(cfg.N1, cfg.N2)
	v := vAt(cfg.CoreRadiusM(), na, cfg.WavelengthM())
	dw := WaveguideDispersion(cfg.N1, RelativeDelta(cfg.N1, cfg.N2), cfg.WavelengthM(), v)
	if dw >= 0 {
		t.Errorf("D_wg = %v, want negative for step-index", dw)
	}
	if math.Abs(dw-(-3.704398)) > 1e-3 {
		t.Errorf("D_wg = %v, want ≈ -3.704398", dw)
	}
}

func TestComposeExampleValues(t *testing.T) {
	cfg := smfCfg()
	res, err := TotalDispersionAt(cfg)
	if err != nil {
		t.Fatalf("TotalDispersionAt: %v", err)
	}
	if math.Abs(res.DMat-3.581818) > 1e-3 {
		t.Errorf("D_mat = %v, want ≈ 3.581818", res.DMat)
	}
	if math.Abs(res.DWg-(-3.704398)) > 1e-3 {
		t.Errorf("D_wg = %v, want ≈ -3.704398", res.DWg)
	}
	if math.Abs(res.DTotal-(-0.122580)) > 1e-3 {
		t.Errorf("D_tot = %v, want ≈ -0.122580", res.DTotal)
	}
}

func TestZeroDispersionWavelengthForSMF(t *testing.T) {
	z, err := ZeroDispersionWavelength(smfCfg(), 1200, 1800)
	if err != nil {
		t.Fatalf("ZeroDispersionWavelength: %v", err)
	}
	if z < 1305 || z > 1318 {
		t.Errorf("zero-D = %v nm, want ≈ 1311.4 nm", z)
	}
}

func TestZeroDispersionNoCrossing(t *testing.T) {
	_, err := ZeroDispersionWavelength(smfCfg(), 1200, 1240)
	if !errors.Is(err, ErrNoZeroDispersion) {
		t.Errorf("want ErrNoZeroDispersion, got %v", err)
	}
}
