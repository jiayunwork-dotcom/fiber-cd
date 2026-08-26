package report

import (
	"fmt"
	"strings"
)

func FmtFloat(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

func FmtSigned(v float64, digits int) string {
	return fmt.Sprintf("%+.*f", digits, v)
}

func FmtPercent(v float64, digits int) string {
	return fmt.Sprintf("%.*f %%", digits, v*100)
}

func FmtNm(v float64, digits int) string {
	return fmt.Sprintf("%.*f nm", digits, v)
}

func FmtVNumber(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

func FmtModes(v float64) string {
	return fmt.Sprintf("%.0f", v)
}

func FmtLambdaRange(lo, hi float64, digits int) string {
	return FmtNm(lo, digits) + " .. " + FmtNm(hi, digits)
}

func FmtRatio(v float64, digits int) string {
	return fmt.Sprintf("%.*f", digits, v)
}

func FmtStepIndex(i, total int) string {
	return fmt.Sprintf("[%d/%d]", i+1, total)
}

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

func PadRight(s string, width int) string {
	if len([]rune(s)) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len([]rune(s)))
}
