package waveguide

import "fiber-cd/internal/model"

type Result struct {
	Config model.Config

	NA                 float64
	Delta              float64
	V                  float64
	CutoffWavelengthNm float64
	Mode               ModeStatus
}
