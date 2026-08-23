package dispersion

// GroupIndex 计算熔融石英在 λ（μm）下的群折射率
//
//	N_g = n − λ·(dn/dλ)
//
// 群折射率决定光脉冲在纤芯材料中的群速度，是色散核算的配套量。
func GroupIndex(lambdaUm float64) float64 {
	s := Silica()
	return s.IndexUm(lambdaUm) - lambdaUm*s.FirstDerivativeUm(lambdaUm)
}

// GroupDelayPerKm 返回每公里群时延 τ = N_g/c，单位 μs/km。
// 换算：τ[s/m] × 1000 m/km × 10⁶ μs/s = τ × 10⁹。
func GroupDelayPerKm(lambdaUm float64) float64 {
	return GroupIndex(lambdaUm) / SpeedOfLightMPS * 1e9
}

// GroupVelocity 返回群速度 v_g = c/N_g，单位 m/s。
func GroupVelocity(lambdaUm float64) float64 {
	return SpeedOfLightMPS / GroupIndex(lambdaUm)
}
