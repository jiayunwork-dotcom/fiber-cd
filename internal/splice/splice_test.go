package splice

import (
	"math"
	"testing"

	"fiber-cd/internal/model"
)

func smf() model.Config {
	return model.Config{N1: 1.4656, N2: 1.4619, CoreDiameterUm: 9, WavelengthNm: 1310}
}

func TestIdenticalFibersZeroOverlapLoss(t *testing.T) {
	loss, err := Evaluate(Pair{A: smf(), B: smf()}, 0)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if math.Abs(loss.Overlap) > 1e-12 {
		t.Errorf("identical overlap loss %v, want 0", loss.Overlap)
	}
	if math.Abs(loss.Offset) > 1e-12 {
		t.Errorf("zero-offset loss %v, want 0", loss.Offset)
	}
	if math.Abs(loss.MFD1-loss.MFD2) > 1e-12 {
		t.Errorf("MFD mismatch on identical fibers: %v vs %v", loss.MFD1, loss.MFD2)
	}
}

func TestOverlapLossSymmetric(t *testing.T) {
	wide := smf()
	wide.CoreDiameterUm = 10
	ab, err := Evaluate(Pair{A: smf(), B: wide}, 0)
	if err != nil {
		t.Fatalf("A→B: %v", err)
	}
	ba, err := Evaluate(Pair{A: wide, B: smf()}, 0)
	if err != nil {
		t.Fatalf("B→A: %v", err)
	}
	if math.Abs(ab.Overlap-ba.Overlap) > 1e-12 {
		t.Errorf("overlap not symmetric: %v vs %v", ab.Overlap, ba.Overlap)
	}
	if ab.Overlap <= 0 {
		t.Errorf("mismatched MFD should produce positive loss, got %v", ab.Overlap)
	}
}

func TestOffsetLossGrowsFasterThanLinear(t *testing.T) {
	d1, err := OffsetLossDB(10, 0.5)
	if err != nil {
		t.Fatalf("0.5: %v", err)
	}
	d2, err := OffsetLossDB(10, 1.0)
	if err != nil {
		t.Fatalf("1.0: %v", err)
	}
	if math.Abs(d2/d1-4) > 1e-9 {
		t.Errorf("doubling offset: loss ratio %v, want 4 (quadratic)", d2/d1)
	}
}

func TestLargerVShrinksMFDAndChangesSplice(t *testing.T) {
	loose := smf()
	tight := smf()
	tight.N2 = 1.450
	a, err := Evaluate(Pair{A: loose, B: loose}, 0)
	if err != nil {
		t.Fatalf("loose: %v", err)
	}
	b, err := Evaluate(Pair{A: tight, B: loose}, 0)
	if err != nil {
		t.Fatalf("tight-loose: %v", err)
	}
	if !(b.MFD1 < a.MFD1) {
		t.Errorf("higher Δ should shrink MFD: tight=%v loose=%v", b.MFD1, a.MFD1)
	}
	if b.Overlap <= 0 {
		t.Errorf("Δ mismatch should give splice loss, got %v", b.Overlap)
	}
}

func TestAngleLossZeroAtNormal(t *testing.T) {
	got, err := AngleLossDB(10, 1.46, 0, 1310)
	if err != nil {
		t.Fatalf("AngleLossDB: %v", err)
	}
	if math.Abs(got) > 1e-12 {
		t.Errorf("zero angle loss %v", got)
	}
	pos, err := AngleLossDB(10, 1.46, 0.01, 1310)
	if err != nil {
		t.Fatalf("0.01 rad: %v", err)
	}
	if pos <= 0 {
		t.Errorf("tilted splice should lose power, got %v", pos)
	}
}

func TestCascadeAdds(t *testing.T) {
	got, err := Cascade([]float64{0.1, 0.2, 0.05})
	if err != nil {
		t.Fatalf("Cascade: %v", err)
	}
	if math.Abs(got-0.35) > 1e-12 {
		t.Errorf("cascade %v, want 0.35", got)
	}
	lin, err := DBToLinear(3)
	if err != nil {
		t.Fatalf("DBToLinear: %v", err)
	}
	back, err := LinearToDB(lin)
	if err != nil {
		t.Fatalf("LinearToDB: %v", err)
	}
	if math.Abs(back-3) > 1e-9 {
		t.Errorf("dB round-trip %v", back)
	}
}

func TestPetermannVsMarcuseSpliceDisagrees(t *testing.T) {
	wide := smf()
	wide.CoreDiameterUm = 10.5
	p := Pair{A: smf(), B: wide}
	delta, err := MFDDefinitionDeltaDB(p, 0.4)
	if err != nil {
		t.Fatalf("MFDDefinitionDeltaDB: %v", err)
	}
	if math.Abs(delta) < 1e-4 {
		t.Errorf("Petermann vs Marcuse splice agreed (%v); definitions must disagree", delta)
	}
	marc, err := EvaluateWithDefinitions(p, 0.4, false)
	if err != nil {
		t.Fatalf("marcuse: %v", err)
	}
	pete, err := EvaluateWithDefinitions(p, 0.4, true)
	if err != nil {
		t.Fatalf("petermann: %v", err)
	}
	if math.Abs(marc.MFD1-pete.MFD1) < 1e-6 {
		t.Errorf("MFD definitions collapsed: %v", marc.MFD1)
	}
	if math.Abs(pete.Total-marc.Total-delta) > 1e-12 {
		t.Errorf("delta %v != pete-marc %v", delta, pete.Total-marc.Total)
	}
}

func TestGapLossZeroAtContact(t *testing.T) {
	z, err := GapLossDB(10, 0, 1310)
	if err != nil {
		t.Fatalf("contact: %v", err)
	}
	if math.Abs(z) > 1e-12 {
		t.Errorf("zero gap loss %v", z)
	}
	g, err := GapLossDB(10, 20, 1310)
	if err != nil {
		t.Fatalf("20um: %v", err)
	}
	if g <= 0 {
		t.Errorf("air gap should lose power, got %v", g)
	}
}
