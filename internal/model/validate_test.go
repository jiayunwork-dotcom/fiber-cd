package model

import (
	"strings"
	"testing"
)

func TestValidateRejectsInvalidFiber(t *testing.T) {
	bad := []Config{
		{N1: 1.4, N2: 1.5, CoreDiameterUm: 9, WavelengthNm: 1310},
		{N1: 1.5, N2: 1.5, CoreDiameterUm: 9, WavelengthNm: 1310},
		{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 0, WavelengthNm: 1310},
		{N1: 1.4656, N2: 1.4619, CoreDiameterUm: -4.5, WavelengthNm: 1310},
		{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9, WavelengthNm: 0},
		{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9, WavelengthNm: -1310},
		{N1: -1.0, N2: 1.5, CoreDiameterUm: 9, WavelengthNm: 1310},
		{N1: 1.4656, N2: 0, CoreDiameterUm: 9, WavelengthNm: 1310},
	}
	for i, cfg := range bad {
		if err := Validate(cfg); err == nil {
			t.Errorf("row %d (%v): expected error, got nil", i, cfg)
		}
	}
}

func TestValidateAcceptsValidFiber(t *testing.T) {
	cfg := Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9.0, WavelengthNm: 1310.0}
	if err := Validate(cfg); err != nil {
		t.Errorf("valid config rejected: %v", err)
	}
}

func TestValidateErrorText(t *testing.T) {
	err := Validate(Config{N1: 1.4, N2: 1.5, CoreDiameterUm: 9, WavelengthNm: 1310})
	if err == nil {
		t.Fatal("expected error for denser cladding, got nil")
	}
	if !strings.Contains(err.Error(), "cladding") {
		t.Errorf("error text %q does not mention %q", err.Error(), "cladding")
	}
}
