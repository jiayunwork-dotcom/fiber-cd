// Package waveguide 实现圆阶跃光纤的导波内核。
//
// 核心物理量：
//
//	NA = √(n1² − n2²)         数值孔径
//	Δ  = (n1 − n2)/n1         相对折射率差
//	V  = 2π·a·NA/λ           归一化频率（a 为纤芯半径，不是直径）
//	单模 iff V < 2.405        截止常数取 J0 首根 2.4048…（圆阶跃钉）
//
// 配套提供：截止波长、模式分类、Gloge 近似 b(V)、Marcuse 光斑尺寸、
// 模式数估计、逆向设计限值（最大单模芯径 / 最大 NA），以及波长扫描
// （sweep）与任意波长探测（probe）。
//
// 依赖方向：waveguide → dispersion（sweep 需要 D_tot），但 dispersion
// 不反向依赖本包，保证导波量可单独测试。
package waveguide
