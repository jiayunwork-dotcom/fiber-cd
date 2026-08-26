package model

func (c Config) RadiusUm() float64 {
	return c.CoreDiameterUm / 2
}

func (c Config) CoreRadiusM() float64 {
	return c.RadiusUm() * 1e-6
}

func (c Config) WavelengthM() float64 {
	return c.WavelengthNm * 1e-9
}

func (c Config) WavelengthUm() float64 {
	return c.WavelengthNm / 1000
}

func (c Config) SweepRange() (startNm, stopNm float64, steps int) {
	startNm = c.SweepStartNm
	if startNm <= 0 {
		startNm = DefaultSweepStartNm
	}
	stopNm = c.SweepStopNm
	if stopNm <= 0 {
		stopNm = DefaultSweepStopNm
	}
	steps = c.SweepSteps
	if steps <= 0 {
		steps = DefaultSweepSteps
	}
	return startNm, stopNm, steps
}
