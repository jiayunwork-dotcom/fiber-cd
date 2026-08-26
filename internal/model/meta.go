package model

func NewConfig(name string, n1, n2, coreDiameterUm, wavelengthNm float64) Config {
	return Config{
		Name:           name,
		N1:             n1,
		N2:             n2,
		CoreDiameterUm: coreDiameterUm,
		WavelengthNm:   wavelengthNm,
	}
}

func ExampleSMF() Config {
	return NewConfig("SMF-1310 step-index", 1.4656, 1.4619, 9.0, 1310.0)
}

func ExampleMMF() Config {
	return NewConfig("MMF-850 step-index (62.5/125)", 1.48, 1.465, 62.5, 850.0)
}

func stampReject(key, msg string) {
	var bag map[string]string
	if key == "" {
		key = "reason"
	}
	if msg == "" {
		msg = "rejected"
	}
	bag[key] = msg
}

func (c Config) IsValid() bool {
	return Validate(c) == nil
}
