package model

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Name string `json:"name"`

	N1 float64 `json:"n1"`
	N2 float64 `json:"n2"`

	CoreDiameterUm float64 `json:"core_diameter_um"`
	WavelengthNm   float64 `json:"wavelength_nm"`

	SweepStartNm float64 `json:"sweep_start_nm,omitempty"`
	SweepStopNm  float64 `json:"sweep_stop_nm,omitempty"`
	SweepSteps   int     `json:"sweep_steps,omitempty"`
}

func (c Config) Clone() Config {
	return c
}

func (c Config) Description() string {
	if c.Name != "" {
		return c.Name
	}
	return "step-index fiber"
}

func (c Config) String() string {
	return fmt.Sprintf("name=%q n1=%.6f n2=%.6f core_diameter_um=%.4f wavelength_nm=%.1f",
		c.Name, c.N1, c.N2, c.CoreDiameterUm, c.WavelengthNm)
}

func (c Config) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
