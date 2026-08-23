// Package dispersion 实现阶跃光纤的色散核算。
//
// 材料色散用钉死的熔融石英（Malitson 1965）Sellmeier 系数对 n(λ)
// 解析求导：
//
//	D_mat = −(λ/c)·d²n/dλ²          [ps/(nm·km)]
//
// 波导色散用阶跃近似公式（Gloge 1971 的 b(V)）：
//
//	D_wg = −(n1·Δ/(λ·c))·V·d²(Vb)/dV²
//
// 总色散 D_tot = D_mat + D_wg，二者在零色散波长附近对消、符号翻转。
// 零色散波长由二分求解，区间内不跨零时明确报错，不编造数值。
//
// 另提供材料斜率 S_mat（闭合形式）、总斜率 S_tot（中心差分）与
// 群折射率 N_g = n − λ·dn/dλ。
//
// 依赖方向：dispersion → model，不依赖 waveguide，避免循环导入。
package dispersion
