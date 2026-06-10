// Package betterminal is a one-import facade over the betterminal toolkit:
// the colors, themes, and grid subpackages. It re-exports the common API so a
// consumer can write everything against a single import:
//
//	import "github.com/GenesisBautista/betterminal"
//
//	t := betterminal.NewGrid(betterminal.Padding(1, 1)).
//		Header("Name", "Score").
//		Row("Ada", 99).
//		Row("Linus", 87)
//	fmt.Print(t)
//
// The subpackages remain importable directly for advanced use; this package
// adds no behavior, only convenient aliases and wrappers.
package betterminal

import (
	"os"

	"github.com/GenesisBautista/betterminal/colors"
	"github.com/GenesisBautista/betterminal/grid"
	"github.com/GenesisBautista/betterminal/themes"
)

// --- colors ---

// Color is an xterm 256-color palette index. See [colors.Color].
type Color = colors.Color

// Style is a terminal text style. See [colors.Style].
type Style = colors.Style

// StyleOption configures a [Style] in [NewStyle].
type StyleOption = colors.Option

// Standard palette colors (indices 0–15).
const (
	Black         = colors.Black
	Red           = colors.Red
	Green         = colors.Green
	Yellow        = colors.Yellow
	Blue          = colors.Blue
	Magenta       = colors.Magenta
	Cyan          = colors.Cyan
	White         = colors.White
	BrightBlack   = colors.BrightBlack
	BrightRed     = colors.BrightRed
	BrightGreen   = colors.BrightGreen
	BrightYellow  = colors.BrightYellow
	BrightBlue    = colors.BrightBlue
	BrightMagenta = colors.BrightMagenta
	BrightCyan    = colors.BrightCyan
	BrightWhite   = colors.BrightWhite
)

// RGB returns a color from the 6×6×6 cube (each component a level 0–5).
func RGB(r, g, b uint8) Color { return colors.RGB(r, g, b) }

// Gray returns a color from the 24-step grayscale ramp (0 dark – 23 light).
func Gray(level uint8) Color { return colors.Gray(level) }

// NewStyle builds a [Style] from options.
func NewStyle(opts ...StyleOption) Style { return colors.New(opts...) }

// Foreground sets the text color.
func Foreground(c Color) StyleOption { return colors.Foreground(c) }

// Background sets the background color.
func Background(c Color) StyleOption { return colors.Background(c) }

// Bold enables bold text.
func Bold() StyleOption { return colors.Bold() }

// Underline enables underlined text.
func Underline() StyleOption { return colors.Underline() }

// IsTerminal reports whether f refers to a terminal.
func IsTerminal(f *os.File) bool { return colors.IsTerminal(f) }

// SetColorEnabled forces ANSI color output on or off (see [colors.Enabled]).
func SetColorEnabled(b bool) { colors.Enabled = b }

// ColorEnabled reports whether ANSI color output is currently enabled.
func ColorEnabled() bool { return colors.Enabled }

// --- themes ---

// Theme maps roles to styles. See [themes.Theme].
type Theme = themes.Theme

// Role names a semantic style slot in a [Theme]. See [themes.Role].
type Role = themes.Role

// Theme roles.
const (
	RoleMessage = themes.Message
	RoleSuccess = themes.Success
	RoleError   = themes.Error
	RoleWarning = themes.Warning
	RoleInfo    = themes.Info
	RoleSystem  = themes.System
	RoleGrid    = themes.Grid
)

// Roles lists the standard roles in display order.
var Roles = themes.Roles

// DefaultTheme returns a theme with a style for every role.
func DefaultTheme() Theme { return themes.Default() }

// LoadTheme reads and parses a theme from a JSON file, reporting any error.
func LoadTheme(path string) (Theme, error) { return themes.Load(path) }

// LoadThemeOrDefault returns the theme at path, or [DefaultTheme] on failure
// (the error is still returned for logging).
func LoadThemeOrDefault(path string) (Theme, error) { return themes.LoadOrDefault(path) }

// --- grid ---

// Table is a console table renderer. See [grid.Table].
type Table = grid.Table

// Cell is a single table value with optional overrides. See [grid.Cell].
type Cell = grid.Cell

// Align controls horizontal placement within a column. See [grid.Align].
type Align = grid.Align

// TableOption configures a [Table] in [NewGrid].
type TableOption = grid.Option

// CellOption configures a single [Cell].
type CellOption = grid.CellOption

// ColumnOption configures a single grid column.
type ColumnOption = grid.ColumnOption

// Column alignment values.
const (
	AlignDefault = grid.AlignDefault
	AlignLeft    = grid.AlignLeft
	AlignRight   = grid.AlignRight
	AlignCenter  = grid.AlignCenter
)

// NewGrid builds a table from options.
func NewGrid(opts ...TableOption) *Table { return grid.New(opts...) }

// C builds a styled cell.
func C(value any, opts ...CellOption) Cell { return grid.C(value, opts...) }

// Padding sets the spaces added inside each cell, left and right.
func Padding(left, right int) TableOption { return grid.Padding(left, right) }

// Alignment sets the table-wide default alignment.
func Alignment(a Align) TableOption { return grid.Alignment(a) }

// FullGrid draws a horizontal border between every data row.
func FullGrid() TableOption { return grid.FullGrid() }

// BorderStyle colors the box-drawing characters.
func BorderStyle(s Style) TableOption { return grid.BorderStyle(s) }

// HeaderStyle colors the header row.
func HeaderStyle(s Style) TableOption { return grid.HeaderStyle(s) }

// AlternatingStyles colors data rows, alternating even/odd.
func AlternatingStyles(even, odd Style) TableOption {
	return grid.AlternatingStyles(even, odd)
}

// CellAlign sets a cell's alignment.
func CellAlign(a Align) CellOption { return grid.CellAlign(a) }

// CellStyle sets a cell's color/attributes.
func CellStyle(s Style) CellOption { return grid.CellStyle(s) }

// ColumnAlign sets a column's default alignment.
func ColumnAlign(a Align) ColumnOption { return grid.ColumnAlign(a) }

// ColumnStyle sets a column's default color/attributes.
func ColumnStyle(s Style) ColumnOption { return grid.ColumnStyle(s) }

// ColumnFormat sets how a column's values are converted to text.
func ColumnFormat(f func(any) string) ColumnOption { return grid.ColumnFormat(f) }
