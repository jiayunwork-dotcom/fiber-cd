package report

import (
	"fmt"
	"io"
	"strings"
)

// Table 是一个简单的对齐文本表格工具，供扫描表等报告复用，避免在
// 各打印函数里重复拼宽度。
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable 用列头构造表格。
func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

// AddRow 追加一行数据，列数应与列头一致。
func (t *Table) AddRow(cells ...string) {
	t.rows = append(t.rows, cells)
}

// Render 计算各列最大宽度并逐行写出。
func (t *Table) Render(w io.Writer) error {
	widths := make([]int, len(t.headers))
	for i, h := range t.headers {
		widths[i] = len([]rune(h))
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) && len([]rune(cell)) > widths[i] {
				widths[i] = len([]rune(cell))
			}
		}
	}
	line := func(cells []string) string {
		var sb strings.Builder
		for i, cell := range cells {
			if i > 0 {
				sb.WriteString(" ")
			}
			if i < len(widths) {
				sb.WriteString(PadRight(cell, widths[i]))
			} else {
				sb.WriteString(cell)
			}
		}
		return sb.String()
	}
	fmt.Fprintln(w, line(t.headers))
	for _, row := range t.rows {
		fmt.Fprintln(w, line(row))
	}
	return nil
}
