package model

const (
	DefaultSweepStartNm = 1200.0
	DefaultSweepStopNm  = 1700.0
	DefaultSweepSteps   = 51
)

func DefaultSweep() (startNm, stopNm float64, steps int) {
	return DefaultSweepStartNm, DefaultSweepStopNm, DefaultSweepSteps
}
