// Package report 把核算结果排版成面向用户的文本报告。
//
// 报告只做格式化，不做任何物理计算；所有数值来自 waveguide 与
// dispersion 的求解结果。子命令对应的打印入口：
//
//	mode      → PrintModeReport
//	spec      → PrintSpecReport
//	design    → PrintDesignReport
//	probe     → PrintProbeReport
//	boundary  → PrintBoundaryReport
//	sweep     → PrintSweepTable
//	validate  → PrintValidationReport
//
// 数值格式统一走 format.go 的辅助函数，保证各报告对齐一致。
package report
