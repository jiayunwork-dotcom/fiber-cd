// Package report 负责把核算结果排版成面向用户的文本报告。这里只做
// 格式化，不做任何物理计算。
package report

import (
	"fmt"
	"strings"
)

// FmtFloat 按固定小数位打印浮点数。
func FmtFloat(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

// FmtSigned 带符号打印浮点数，用于色散值（负号不能丢）。
func FmtSigned(v float64, digits int) string {
	return fmt.Sprintf("%+.*f", digits, v)
}

// FmtPercent 把相对折射率差（0–1）打印成百分数。
func FmtPercent(v float64, digits int) string {
	return fmt.Sprintf("%.*f %%", digits, v*100)
}

// FmtNm 打印带单位的波长（nm）。
func FmtNm(v float64, digits int) string {
	return fmt.Sprintf("%.*f nm", digits, v)
}

// FmtVNumber 打印归一化频率 V。
func FmtVNumber(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

// FmtModes 打印模式数估计（保留 0 位小数）。
func FmtModes(v float64) string {
	return fmt.Sprintf("%.0f", v)
}

// FmtLambdaRange 打印一段波长区间 "lo nm .. hi nm"。
func FmtLambdaRange(lo, hi float64, digits int) string {
	return FmtNm(lo, digits) + " .. " + FmtNm(hi, digits)
}

// FmtRatio 打印无单位比值，如芯径相对最大值的百分比。
func FmtRatio(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

// FmtStepIndex 打印档位编号，如 "[12/51]"，用于长扫描表的分页行首。
func FmtStepIndex(i, total int) string {
	return fmt.Sprintf("[%d/%d]", i+1, total)
}

// Truncate 把字符串截断到 max 个字符，超长用省略号收尾。
func Truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

// PadRight 右补空格到固定宽度，供扫描表列头对齐。
func PadRight(s string, width int) string {
	if len([]rune(s)) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len([]rune(s)))
}
