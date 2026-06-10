// Package grid renders aligned, bordered tables to the terminal. It is an
// idiomatic Go take on the author's C# ConsoleTableMaker.
//
//	t := grid.New(grid.Padding(1, 1)).
//		Header("Name", "Score").
//		Row("Ada", 99).
//		Row("Linus", 87)
//	fmt.Print(t)
//
// Styling reuses the colors package. Alignment and color resolve with a
// cell › column › table precedence, so a cell override beats a column
// override beats the table default.
package grid

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/GenesisBautista/betterminal/colors"
)

// Align controls horizontal text placement within a column. The zero value,
// AlignDefault, inherits from the column, then the table (which defaults to
// AlignLeft).
type Align uint8

const (
	AlignDefault Align = iota
	AlignLeft
	AlignRight
	AlignCenter
)

// Cell is a single table value with optional per-cell overrides. Build one
// with [C]; plain values passed to [Table.Row] are wrapped automatically.
type Cell struct {
	Value    any
	Align    Align
	Style    colors.Style
	hasStyle bool
}

// CellOption overrides styling or alignment for a single [Cell].
type CellOption func(*Cell)

// C builds a styled cell, e.g. C("total", grid.CellAlign(grid.AlignRight)).
func C(value any, opts ...CellOption) Cell {
	c := Cell{Value: value}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// CellAlign sets a cell's alignment.
func CellAlign(a Align) CellOption { return func(c *Cell) { c.Align = a } }

// CellStyle sets a cell's color/attributes.
func CellStyle(s colors.Style) CellOption {
	return func(c *Cell) { c.Style, c.hasStyle = s, true }
}

type columnOpts struct {
	align    Align
	style    colors.Style
	hasStyle bool
	format   func(any) string
}

// Table accumulates header and data rows and renders them as a box-drawn
// grid. The zero value is not usable; start with [New].
type Table struct {
	headers    []Cell
	hasHeaders bool
	rows       [][]Cell

	padLeft, padRight int
	align             Align
	fullGrid          bool

	alternating bool
	dataStyles  [2]colors.Style
	borderStyle colors.Style
	headerStyle colors.Style

	columns map[int]columnOpts

	cols int   // first-seen column count, for validation
	err  error // first build error
}

// Option configures a [Table] in [New].
type Option func(*Table)

// New builds a table from options.
func New(opts ...Option) *Table {
	t := &Table{columns: map[int]columnOpts{}}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Padding sets the number of spaces added inside each cell, left and right.
func Padding(left, right int) Option {
	return func(t *Table) { t.padLeft, t.padRight = left, right }
}

// Alignment sets the table-wide default alignment.
func Alignment(a Align) Option { return func(t *Table) { t.align = a } }

// FullGrid draws a horizontal border between every data row, not just under
// the header.
func FullGrid() Option { return func(t *Table) { t.fullGrid = true } }

// BorderStyle colors the box-drawing characters.
func BorderStyle(s colors.Style) Option { return func(t *Table) { t.borderStyle = s } }

// HeaderStyle colors the header row. It always wins over column and cell
// styles for header cells.
func HeaderStyle(s colors.Style) Option { return func(t *Table) { t.headerStyle = s } }

// AlternatingStyles colors data rows, alternating between even and odd rows.
// Cell and column styles still take precedence.
func AlternatingStyles(even, odd colors.Style) Option {
	return func(t *Table) { t.alternating, t.dataStyles = true, [2]colors.Style{even, odd} }
}

// Header sets the header row. Values may be plain or [Cell] values.
func (t *Table) Header(values ...any) *Table {
	t.headers = toCells(values)
	t.hasHeaders = true
	t.checkCols(len(t.headers))
	return t
}

// Row appends a data row. Values may be plain or [Cell] values. Returns the
// table for chaining.
func (t *Table) Row(values ...any) *Table {
	cells := toCells(values)
	t.rows = append(t.rows, cells)
	t.checkCols(len(cells))
	return t
}

// Column applies options (alignment, style, formatter) to the column at index i.
func (t *Table) Column(i int, opts ...ColumnOption) *Table {
	co := t.columns[i]
	for _, opt := range opts {
		opt(&co)
	}
	t.columns[i] = co
	return t
}

// ColumnOption configures a single column.
type ColumnOption func(*columnOpts)

// ColumnAlign sets a column's default alignment.
func ColumnAlign(a Align) ColumnOption { return func(c *columnOpts) { c.align = a } }

// ColumnStyle sets a column's default color/attributes.
func ColumnStyle(s colors.Style) ColumnOption {
	return func(c *columnOpts) { c.style, c.hasStyle = s, true }
}

// ColumnFormat sets how a column's values are converted to text. The default
// is fmt.Sprint.
func ColumnFormat(f func(any) string) ColumnOption {
	return func(c *columnOpts) { c.format = f }
}

// Err reports the first build error, such as a row whose column count does
// not match the rest of the table. Rendering still proceeds best-effort.
func (t *Table) Err() error { return t.err }

func toCells(values []any) []Cell {
	cells := make([]Cell, len(values))
	for i, v := range values {
		if c, ok := v.(Cell); ok {
			cells[i] = c
		} else {
			cells[i] = Cell{Value: v}
		}
	}
	return cells
}

func (t *Table) checkCols(n int) {
	if t.cols == 0 {
		t.cols = n
		return
	}
	if n != t.cols && t.err == nil {
		t.err = fmt.Errorf("grid: row has %d columns, want %d", n, t.cols)
	}
}

// --- rendering ---

// String renders the table to a string. It implements [fmt.Stringer], so a
// *Table can be passed straight to fmt.Print.
func (t *Table) String() string {
	ncols := t.columnCount()
	if ncols == 0 {
		return ""
	}

	text, styles, aligns := t.layout(ncols)
	widths := columnWidths(text, ncols)

	var b strings.Builder
	t.writeBorder(&b, "┌", "┬", "┐", widths)

	line := 0
	if t.hasHeaders {
		t.writeRow(&b, text[line], styles[line], aligns[line], widths)
		t.writeBorder(&b, "├", "┼", "┤", widths)
		line++
	}
	for r := 0; line < len(text); r, line = r+1, line+1 {
		t.writeRow(&b, text[line], styles[line], aligns[line], widths)
		if t.fullGrid && r < len(t.rows)-1 {
			t.writeBorder(&b, "├", "┼", "┤", widths)
		}
	}

	t.writeBorder(&b, "└", "┴", "┘", widths)
	return b.String()
}

// Render writes the table to w. It returns any write error, or the build
// error from [Table.Err] if the table is malformed.
func (t *Table) Render(w io.Writer) error {
	if _, err := io.WriteString(w, t.String()); err != nil {
		return err
	}
	return t.err
}

func (t *Table) columnCount() int {
	n := 0
	if t.hasHeaders {
		n = len(t.headers)
	}
	for _, row := range t.rows {
		if len(row) > n {
			n = len(row)
		}
	}
	return n
}

// layout produces the formatted text, resolved style, and resolved alignment
// for every cell, header row first.
func (t *Table) layout(ncols int) (text [][]string, styles [][]colors.Style, aligns [][]Align) {
	add := func(row []Cell, header bool, dataRow int) {
		tline := make([]string, ncols)
		sline := make([]colors.Style, ncols)
		aline := make([]Align, ncols)
		for c := 0; c < ncols; c++ {
			var cell Cell
			if c < len(row) {
				cell = row[c]
			}
			tline[c] = t.format(c, header, cell.Value)
			sline[c] = t.cellStyle(header, dataRow, c, cell)
			aline[c] = t.cellAlign(c, cell)
		}
		text = append(text, tline)
		styles = append(styles, sline)
		aligns = append(aligns, aline)
	}

	if t.hasHeaders {
		add(t.headers, true, 0)
	}
	for i, row := range t.rows {
		add(row, false, i)
	}
	return text, styles, aligns
}

// format converts a cell value to text. The per-column formatter applies to
// data cells only; header labels are always rendered with fmt.Sprint.
func (t *Table) format(col int, header bool, v any) string {
	if v == nil {
		v = ""
	}
	if !header {
		if co, ok := t.columns[col]; ok && co.format != nil {
			return co.format(v)
		}
	}
	return fmt.Sprint(v)
}

func (t *Table) cellStyle(header bool, dataRow, col int, cell Cell) colors.Style {
	if header {
		return t.headerStyle
	}
	if cell.hasStyle {
		return cell.Style
	}
	if co, ok := t.columns[col]; ok && co.hasStyle {
		return co.style
	}
	if t.alternating {
		return t.dataStyles[dataRow%2]
	}
	return colors.Style{}
}

func (t *Table) cellAlign(col int, cell Cell) Align {
	if cell.Align != AlignDefault {
		return cell.Align
	}
	if co, ok := t.columns[col]; ok && co.align != AlignDefault {
		return co.align
	}
	if t.align != AlignDefault {
		return t.align
	}
	return AlignLeft
}

func columnWidths(text [][]string, ncols int) []int {
	widths := make([]int, ncols)
	for _, line := range text {
		for c, s := range line {
			if w := utf8.RuneCountInString(s); w > widths[c] {
				widths[c] = w
			}
		}
	}
	return widths
}

func (t *Table) writeBorder(b *strings.Builder, left, mid, right string, widths []int) {
	var line strings.Builder
	line.WriteString(left)
	for i, w := range widths {
		if i > 0 {
			line.WriteString(mid)
		}
		line.WriteString(strings.Repeat("─", t.padLeft+w+t.padRight))
	}
	line.WriteString(right)
	b.WriteString(t.borderStyle.Render(line.String()))
	b.WriteByte('\n')
}

func (t *Table) writeRow(b *strings.Builder, text []string, styles []colors.Style, aligns []Align, widths []int) {
	bar := t.borderStyle.Render("│")
	pad := strings.Repeat(" ", t.padLeft)
	rpad := strings.Repeat(" ", t.padRight)

	b.WriteString(bar)
	for c, s := range text {
		content := pad + alignWithin(s, widths[c], aligns[c]) + rpad
		b.WriteString(styles[c].Render(content))
		b.WriteString(bar)
	}
	b.WriteByte('\n')
}

func alignWithin(s string, width int, a Align) string {
	gap := width - utf8.RuneCountInString(s)
	if gap <= 0 {
		return s
	}
	switch a {
	case AlignRight:
		return strings.Repeat(" ", gap) + s
	case AlignCenter:
		l := gap / 2
		return strings.Repeat(" ", l) + s + strings.Repeat(" ", gap-l)
	default:
		return s + strings.Repeat(" ", gap)
	}
}
