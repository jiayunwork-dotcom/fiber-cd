package splice

import (
	"fmt"
	"math"
)

func Cascade(losses []float64) (float64, error) {
	sum := 0.0
	for _, l := range losses {
		if l < 0 {
			return 0, fmt.Errorf("splice: cascaded loss must be >= 0")
		}
		sum += l
	}
	return sum, nil
}

func MeanMFD(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("splice: empty MFD list")
	}
	s := 0.0
	for _, v := range values {
		if v <= 0 {
			return 0, fmt.Errorf("splice: MFD must be > 0")
		}
		s += v
	}
	return s / float64(len(values)), nil
}

func RMSOffset(offsets []float64) (float64, error) {
	if len(offsets) == 0 {
		return 0, fmt.Errorf("splice: empty offset list")
	}
	s := 0.0
	for _, v := range offsets {
		if v < 0 {
			return 0, fmt.Errorf("splice: offset must be >= 0")
		}
		s += v * v
	}
	return math.Sqrt(s / float64(len(offsets))), nil
}
