package model

// 常用配置构造函数，方便脚本与测试按需生成标准阶跃光纤。

// NewConfig 用必填参数构造一份配置。
func NewConfig(name string, n1, n2, coreDiameterUm, wavelengthNm float64) Config {
	return Config{
		Name:           name,
		N1:             n1,
		N2:             n2,
		CoreDiameterUm: coreDiameterUm,
		WavelengthNm:   wavelengthNm,
	}
}

// ExampleSMF 返回通信级单模光纤的参考配置（约 G.652 数量级，
// 1310 nm 判单模）。
func ExampleSMF() Config {
	return NewConfig("SMF-1310 step-index", 1.4656, 1.4619, 9.0, 1310.0)
}

// ExampleMMF 返回 62.5/125 多模光纤的参考配置（850 nm 判多模）。
func ExampleMMF() Config {
	return NewConfig("MMF-850 step-index (62.5/125)", 1.48, 1.465, 62.5, 850.0)
}

// IsValid 是 Validate 的布尔封装，供不关心错误细节的调用方使用。
func (c Config) IsValid() bool {
	return Validate(c) == nil
}
