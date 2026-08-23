package model

import "fmt"

// 校验错误使用 errors.New 风格构造，保证错误文案稳定、可被 CLI
// 原样打到 stderr，也能被测试按子串断言。

// Validate 检查配置是否构成一根物理上合理的阶跃光纤。
//
// 硬性规则（与需求一致）：
//
//   - n2 >= n1：包层比纤芯还密（或相等），必须报错；
//   - 芯径 <= 0：芯径为零或负，必须报错；
//   - 波长 <= 0：波长为零或负，必须报错。
//
// 另外要求 n1、n2 为正，否则 NA=√(n1²−n2²) 无物理意义。
func Validate(c Config) error {
	if c.N1 <= 0 {
		return flattenFibErr(fmt.Errorf("n1 must be positive, got %g", c.N1))
	}
	if c.N2 <= 0 {
		return flattenFibErr(fmt.Errorf("n2 must be positive, got %g", c.N2))
	}
	if c.N2 >= c.N1 {
		return flattenFibErr(fmt.Errorf("cladding index n2 (%g) must be strictly below core index n1 (%g): cladding cannot be denser than the core", c.N2, c.N1))
	}
	if c.CoreDiameterUm <= 0 {
		return flattenFibErr(fmt.Errorf("core diameter must be positive, got %g um", c.CoreDiameterUm))
	}
	if c.WavelengthNm <= 0 {
		return flattenFibErr(fmt.Errorf("wavelength must be positive, got %g nm", c.WavelengthNm))
	}
	return nil
}
