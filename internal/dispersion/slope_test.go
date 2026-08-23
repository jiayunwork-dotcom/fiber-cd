package dispersion

import (
	"math"
	"testing"
)

func TestGroupIndexRange(t *testing.T) {
	ng := GroupIndex(1.31)
	if ng < 1.46 || ng > 1.47 {
		t.Errorf("group index at 1310 nm = %v, want ≈ 1.4616", ng)
	}
	if math.Abs(ng-1.461640) > 1e-4 {
		t.Errorf("group index = %v, want ≈ 1.461640", ng)
	}
	delay := GroupDelayPerKm(1.31)
	if delay < 4.8 || delay > 5.0 {
		t.Errorf("group delay = %v us/km, want ≈ 4.875 us/km", delay)
	}
}

func TestMaterialDispersionSlopeSign(t *testing.T) {
	s := MaterialDispersionSlopeUm(1.31)
	if s <= 0 {
		t.Errorf("S_mat(1310nm) = %v, want positive", s)
	}
	if math.Abs(s-0.0925) > 0.01 {
		t.Errorf("S_mat(1310nm) = %v, want ≈ 0.0925", s)
	}
}

func TestTotalDispersionSlopeNearZero(t *testing.T) {
	cfg := smfCfg()
	s, err := TotalDispersionSlope(cfg, 1310.0)
	if err != nil {
		t.Fatalf("TotalDispersionSlope: %v", err)
	}
	if math.Abs(s-0.0897) > 0.005 {
		t.Errorf("S_tot(1310nm) = %v, want ≈ 0.0897", s)
	}
}

func TestSlopeAtOperatingComposition(t *testing.T) {
	cfg := smfCfg()
	res, err := SlopeAtOperating(cfg)
	if err != nil {
		t.Fatalf("SlopeAtOperating: %v", err)
	}
	want := res.SMaterial + res.SWaveguide
	if math.Abs(res.STotal-want) > 1e-9 {
		t.Errorf("S_tot = %v, want S_mat + S_wg = %v", res.STotal, want)
	}
}

// n”'(λ) ≈ (n”(λ+h) − n”(λ−h)) / 2h。
func TestThirdDerivativeFiniteDifference(t *testing.T) {
	s := Silica()
	lambda := 1.31
	h := 1e-4
	n2p := s.SecondDerivativeUm(lambda + h)
	n2m := s.SecondDerivativeUm(lambda - h)
	fd := (n2p - n2m) / (2 * h)
	closed := s.ThirdDerivativeUm(lambda)
	if math.Abs(fd-closed) > 1e-2*math.Abs(closed) {
		t.Errorf("third derivative mismatch: finite-diff %v vs closed %v", fd, closed)
	}
}
