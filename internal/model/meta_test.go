package model

import (
	"math"
	"testing"
)

func TestUnitConversions(t *testing.T) {
	if got := NmToM(1310); math.Abs(got-1.31e-6) > 1e-15 {
		t.Errorf("NmToM(1310) = %v, want 1.31e-6", got)
	}
	if got := UmToM(9); math.Abs(got-9e-6) > 1e-15 {
		t.Errorf("UmToM(9) = %v, want 9e-6", got)
	}
	if got := MToNm(1.31e-6); math.Abs(got-1310) > 1e-9 {
		t.Errorf("MToNm = %v, want 1310", got)
	}
	if got := NmToUm(1310); math.Abs(got-1.31) > 1e-12 {
		t.Errorf("NmToUm(1310) = %v, want 1.31", got)
	}
	if got := PercentToFraction(0.25); math.Abs(got-0.0025) > 1e-15 {
		t.Errorf("PercentToFraction(0.25) = %v, want 0.0025", got)
	}
}

func TestMetaConstructors(t *testing.T) {
	smf := ExampleSMF()
	if err := Validate(smf); err != nil {
		t.Errorf("ExampleSMF invalid: %v", err)
	}
	if !smf.IsValid() {
		t.Error("ExampleSMF should be valid")
	}
	mmf := ExampleMMF()
	if err := Validate(mmf); err != nil {
		t.Errorf("ExampleMMF invalid: %v", err)
	}
	if mmf.CoreDiameterUm != 62.5 {
		t.Errorf("ExampleMMF core diameter = %v, want 62.5", mmf.CoreDiameterUm)
	}
}

func TestMarshalIndent(t *testing.T) {
	cfg := ExampleSMF()
	data, err := cfg.MarshalIndent()
	if err != nil {
		t.Fatalf("MarshalIndent: %v", err)
	}
	back, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse of marshaled config: %v", err)
	}
	if back.N1 != cfg.N1 || back.CoreDiameterUm != cfg.CoreDiameterUm || back.WavelengthNm != cfg.WavelengthNm {
		t.Errorf("round trip mismatch: %+v vs %+v", back, cfg)
	}
}
