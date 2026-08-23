package dispersion

// FirstDerivativeUm 解析计算 dn/dλ，λ 以 μm 计。对
//
//	n² = 1 + Σ B_i λ²/(λ² − C_i)
//
// 两边求导可得 n' = −(λ/n)·Σ B_i·C_i/(λ² − C_i)²。
func (s Sellmeier) FirstDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		denom := lambdaUm*lambdaUm - t.C
		sum += t.B * t.C / (denom * denom)
	}
	return -lambdaUm / n * sum
}

// SecondDerivativeUm 解析计算 d²n/dλ²，λ 以 μm 计。对
//
//	2n·n' = Σ −2λ·C_i·B_i/(λ² − C_i)²
//
// 再求一次导并整理：
//
//	n'' = [ Σ C_i·B_i·(3λ² + C_i)/(λ² − C_i)³ − (n')² ] / n
func (s Sellmeier) SecondDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	n1 := s.FirstDerivativeUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		denom := l2 - t.C
		cube := denom * denom * denom
		sum += t.C * t.B * (3*l2 + t.C) / cube
	}
	return (sum - n1*n1) / n
}

// ThirdDerivativeUm 解析计算 d³n/dλ³，λ 以 μm 计，用于材料色散斜率
// S_mat = −(1/c)·(n” + λ·n”') 的闭合形式。由
//
//	2(n')² + 2n·n'' = Σ 2·C·B·(3λ²+C)/(λ²−C)³
//
// 再求一次导：
//
//	n''' = ( −24λ·Σ C·B·(λ²+C)/(λ²−C)⁴ − 4n'n'' − 2(n'')² ) / (2n)
func (s Sellmeier) ThirdDerivativeUm(lambdaUm float64) float64 {
	n := s.IndexUm(lambdaUm)
	n1 := s.FirstDerivativeUm(lambdaUm)
	n2 := s.SecondDerivativeUm(lambdaUm)
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		denom := l2 - t.C
		pow4 := denom * denom * denom * denom
		sum += t.C * t.B * (l2 + t.C) / pow4
	}
	return (-24*lambdaUm*sum - 4*n1*n2 - 2*n2*n2) / (2 * n)
}
