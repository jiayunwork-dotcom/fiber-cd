// Package dispersion 实现阶跃光纤的色散核算：材料色散来自熔融石英
// （Malitson 1965）Sellmeier 系数对 n(λ) 的解析求导，波导色散采用
// 阶跃近似（Gloge 1971）的钉文献公式，总色散 D_tot = D_mat + D_wg。
package dispersion

import "math"

// Term 是 Sellmeier 方程的一个谐振项：
//
//	n²(λ) = 1 + Σ B_i·λ² / (λ² − C_i)
//
// B_i 无量纲，C_i 为谐振波长的平方，单位 μm²。
type Term struct {
	B float64
	C float64
}

// Sellmeier 用一组谐振项描述某种玻璃材料在 λ（μm）下的折射率。
type Sellmeier struct {
	Name  string
	Terms []Term
}

// Silica 返回熔融石英（fused silica, SiO₂）的 Malitson 1965 系数，
// 这是本工具钉死的材料色散基准，不随配置改变。
func Silica() Sellmeier {
	return Sellmeier{
		Name: "fused silica (Malitson 1965)",
		Terms: []Term{
			{B: 0.6961663, C: 0.00467914826},
			{B: 0.4079426, C: 0.0135120631},
			{B: 0.8974794, C: 97.9340025},
		},
	}
}

// IndexUm 计算 λ（μm）下的折射率 n(λ)。
func (s Sellmeier) IndexUm(lambdaUm float64) float64 {
	sum := 0.0
	for _, t := range s.Terms {
		l2 := lambdaUm * lambdaUm
		sum += t.B * l2 / (l2 - t.C)
	}
	return math.Sqrt(1 + sum)
}
