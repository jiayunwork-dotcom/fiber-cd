package model

// RadiusUm 返回纤芯半径。用户输入的是芯径（直径），V 数公式里的
// a 是半径，这里完成 芯径 → 半径 的唯一折算点，禁止在别处再除以 2。
func (c Config) RadiusUm() float64 {
	return c.CoreDiameterUm / 2
}

// CoreRadiusM 返回以米为单位的纤芯半径，供 V 数、截止波长等 SI 公式使用。
func (c Config) CoreRadiusM() float64 {
	return c.RadiusUm() * 1e-6
}

// WavelengthM 返回以米为单位的工作波长。
func (c Config) WavelengthM() float64 {
	return c.WavelengthNm * 1e-9
}

// WavelengthUm 返回以 μm 为单位的工作波长，供 Sellmeier 求值。
func (c Config) WavelengthUm() float64 {
	return c.WavelengthNm / 1000
}

// SweepRange 决定波长扫描区间与步数，优先级：配置字段 > 包级默认值。
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
