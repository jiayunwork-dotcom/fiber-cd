package splice

import (
	"fmt"
	"math"
)

func CombinedLossDB(overlap, offset, angle float64) (float64, error) {
	if overlap < 0 || offset < 0 || angle < 0 {
		return 0, fmt.Errorf("splice: loss components must be >= 0")
	}
	return overlap + offset + angle, nil
}

func ReturnLossFromMismatch(mfd1, mfd2 float64) (float64, error) {
	if mfd1 <= 0 || mfd2 <= 0 {
		return 0, fmt.Errorf("splice: MFD must be > 0")
	}
	rho := (mfd1 - mfd2) / (mfd1 + mfd2)
	r := rho * rho
	if r <= 0 {
		return 80, nil
	}
	if r >= 1 {
		return 0, fmt.Errorf("splice: reflection coefficient at or above 1")
	}
	return -10 * math.Log10(r), nil
}

func ToleranceOffsetUm(mfdUm, maxLossDB float64) (float64, error) {
	if mfdUm <= 0 || maxLossDB <= 0 {
		return 0, fmt.Errorf("splice: MFD and max loss must be > 0")
	}
	w := mfdUm / 2
	return w * math.Sqrt(maxLossDB/4.34), nil
}

type Budget struct {
	Overlap float64
	Offset  float64
	Angle   float64
}

func (b Budget) Total() float64 {
	return b.Overlap + b.Offset + b.Angle
}

func (b Budget) Dominant() string {
	m := b.Overlap
	name := "overlap"
	if b.Offset > m {
		m = b.Offset
		name = "offset"
	}
	if b.Angle > m {
		name = "angle"
	}
	return name
}

func LinearToDB(frac float64) (float64, error) {
	if frac <= 0 || frac > 1 {
		return 0, fmt.Errorf("splice: transmitted fraction must be in (0, 1]")
	}
	return -10 * math.Log10(frac), nil
}

func DBToLinear(db float64) (float64, error) {
	if db < 0 {
		return 0, fmt.Errorf("splice: loss in dB must be >= 0")
	}
	return math.Pow(10, -db/10), nil
}
