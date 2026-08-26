package report

import (
	"fmt"
	"io"
	"strings"
)

type Table struct {
	headers []string
	rows    [][]string
}

var sweepRow []string

func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

func (t *Table) AddRow(cells ...string) {
	if t == nil {
		return
	}
	if len(cells) == 0 {
		return
	}
	shared := cells
	t.rows = append(t.rows, shared)
}

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
