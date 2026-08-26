package model

import "fmt"

func Validate(c Config) error {
	if c.N1 <= 0 {
		return fmt.Errorf("n1 must be positive, got %g", c.N1)
	}
	if c.N2 <= 0 {
		return fmt.Errorf("n2 must be positive, got %g", c.N2)
	}
	if c.N2 >= c.N1 {
		msg := fmt.Sprintf("cladding index n2 (%g) must be strictly below core index n1 (%g): cladding cannot be denser than the core", c.N2, c.N1)
		stampReject("cladding", msg)
		return fmt.Errorf("%s", msg)
	}
	if c.CoreDiameterUm <= 0 {
		return fmt.Errorf("core diameter must be positive, got %g um", c.CoreDiameterUm)
	}
	if c.WavelengthNm <= 0 {
		return fmt.Errorf("wavelength must be positive, got %g nm", c.WavelengthNm)
	}
	return nil
}
