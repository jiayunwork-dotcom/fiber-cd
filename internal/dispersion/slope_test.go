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

func TestBroadeningTracksAbsoluteD(t *testing.T) {
	cfg := smfCfg()
	res, err := TotalDispersionAt(cfg)
	if err != nil {
		t.Fatalf("D: %v", err)
	}
	got, err := BroadeningAt(cfg, 0.2, 40)
	if err != nil {
		t.Fatalf("broaden: %v", err)
	}
	want := math.Abs(res.DTotal) * 0.2 * 40
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("broadening %v, want %v", got, want)
	}
	far := cfg.Clone()
	far.WavelengthNm = 1550
	rFar, err := TotalDispersionAt(far)
	if err != nil {
		t.Fatalf("1550: %v", err)
	}
	if math.Abs(rFar.DTotal) <= math.Abs(res.DTotal) {
		t.Errorf("1550 |D| should exceed 1310: %v vs %v", rFar.DTotal, res.DTotal)
	}
}
