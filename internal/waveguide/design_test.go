package waveguide

import (
	"math"
	"testing"
)

func TestSpotSizeRange(t *testing.T) {
	cfg := smfCfg()
	res, err := Analyze(cfg)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	w0 := SpotRadiusUm(cfg.RadiusUm(), res.V)
	if w0 <= cfg.RadiusUm()*0.6 || w0 >= cfg.RadiusUm()*1.5 {
		t.Errorf("spot radius w0 = %v um out of plausible band for a=%v um", w0, cfg.RadiusUm())
	}
	mfd := ModeFieldDiameterUm(cfg.RadiusUm(), res.V)
	if math.Abs(mfd-10.380) > 0.05 {
		t.Errorf("MFD = %v um, want ≈ 10.380 um", mfd)
	}
}

func TestEffectiveAreaGrowsWithV(t *testing.T) {
	base := smfCfg()
	big := base.Clone()
	big.CoreDiameterUm = base.CoreDiameterUm * 2

	na := NumericalAperture(base.N1, base.N2)
	v1 := VNumber(base.CoreRadiusM(), na, base.WavelengthM())
	v2 := VNumber(big.CoreRadiusM(), na, big.WavelengthM())

	a1 := EffectiveAreaUm2(base.RadiusUm(), v1)
	a2 := EffectiveAreaUm2(big.RadiusUm(), v2)
	if a2 <= a1 {
		t.Errorf("effective area should grow with core diameter: %v vs %v", a2, a1)
	}
}

func TestEffectiveModeCount(t *testing.T) {
	if got := EffectiveModeCount(2.0); got != 1 {
		t.Errorf("EffectiveModeCount(2.0) = %v, want 1 (single-mode)", got)
	}
	if got := TotalModeCountEstimate(48.5); got < 1000 || got > 1300 {
		t.Errorf("TotalModeCountEstimate(48.5) = %v, want ≈ 1176", got)
	}
	if math.Abs(TotalModeCountEstimate(48.5)-1176.1) > 1.0 {
		t.Errorf("TotalModeCountEstimate(48.5) = %v, want ≈ 1176.1", TotalModeCountEstimate(48.5))
	}
}

func TestMaxSingleModeDiameter(t *testing.T) {
	cfg := smfCfg()
	na := NumericalAperture(cfg.N1, cfg.N2)
	maxD := MaxSingleModeDiameterUm(na, cfg.WavelengthNm)
	if math.Abs(maxD-9.636) > 0.02 {
		t.Errorf("max single-mode diameter = %v um, want ≈ 9.636 um", maxD)
	}
	if cfg.CoreDiameterUm >= maxD {
		t.Errorf("example diameter %v um must be below max %v um", cfg.CoreDiameterUm, maxD)
	}
	if 18.0 <= maxD {
		t.Error("doubled diameter must exceed the single-mode limit")
	}

	maxNA := MaxSingleModeNA(cfg.CoreDiameterUm, cfg.WavelengthNm)
	if math.Abs(maxNA-0.111428) > 1e-4 {
		t.Errorf("max NA = %v, want ≈ 0.111428", maxNA)
	}
	if na >= maxNA {
		t.Errorf("example NA %v must be below max %v", na, maxNA)
	}
}

func TestSingleModeMargin(t *testing.T) {
	if got := SingleModeMargin(2.2463); got <= 0 {
		t.Errorf("margin for V=2.2463 = %v, want positive", got)
	}
	if got := SingleModeMargin(4.49); got >= 0 {
		t.Errorf("margin for V=4.49 = %v, want negative", got)
	}
}
