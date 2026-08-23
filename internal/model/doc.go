// Package model 定义阶跃光纤的输入配置 schema、JSON 装载、非法输入
// 校验与单位换算。
//
// 配置以 JSON 文件给出：纤芯/包层折射率 n1、n2，芯径
// core_diameter_um（直径，单位 μm）与工作波长 wavelength_nm（单位
// nm）。硬性校验：n2 ≥ n1 报错（包层比纤芯还密）、芯径 ≤ 0 报错、
// 波长 ≤ 0 报错。芯径到半径的折算只在这里完成一次。
//
// 本包不依赖任何计算逻辑，是 waveguide / dispersion / report 三层
// 共用的数据底座。
package model
