// Package colors builds ANSI escape sequences for the xterm 256-color
// palette and applies them to strings.
//
// A [Color] is a palette index. A [Style] couples an optional foreground
// and background color with text attributes (bold, underline) and is built
// with functional options:
//
//	s := colors.New(colors.Foreground(colors.RGB(0, 5, 0)), colors.Bold())
//	fmt.Println(s.Render("ok"))
//
// Output is gated by the package-level [Enabled] flag, which defaults to
// false when the NO_COLOR environment variable is set (see https://no-color.org).
package colors

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Color is an index into the xterm 256-color palette:
//
//	0–15    standard and bright terminal colors
//	16–231  the 6×6×6 RGB cube
//	232–255 the 24-step grayscale ramp
type Color uint8

// Standard palette colors (indices 0–15).
const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// RGB returns the palette color nearest the given point in the xterm 6×6×6
// color cube. Each component is a level in the range 0–5; larger values are
// clamped to 5.
func RGB(r, g, b uint8) Color {
	return Color(16 + 36*clampLevel(r) + 6*clampLevel(g) + clampLevel(b))
}

// Gray returns a color from the 24-step grayscale ramp, where level 0 is
// darkest and 23 is lightest. Larger values are clamped to 23.
func Gray(level uint8) Color {
	if level > 23 {
		level = 23
	}
	return Color(232 + level)
}

func clampLevel(v uint8) uint8 {
	if v > 5 {
		return 5
	}
	return v
}

// Enabled controls whether [Style.Render] emits ANSI escape codes. It
// defaults to false when NO_COLOR is set in the environment, true otherwise.
// Set it explicitly to force color on or off, e.g.:
//
//	colors.Enabled = colors.IsTerminal(os.Stdout)
var Enabled = func() bool {
	_, noColor := os.LookupEnv("NO_COLOR")
	return !noColor
}()

// IsTerminal reports whether f refers to a terminal (a character device),
// which is a reasonable signal for whether to emit color.
func IsTerminal(f *os.File) bool {
	fi, err := f.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

const (
	csi   = "\x1b["
	reset = "\x1b[0m"
)

// Style is an immutable terminal text style. The zero Style renders text
// unchanged. Build non-trivial styles with [New] and the With/Foreground/
// Background/Bold/Underline options.
type Style struct {
	fg, bg    Color
	hasFg     bool
	hasBg     bool
	bold      bool
	underline bool
}

// Option configures a [Style] in [New].
type Option func(*Style)

// New builds a Style from zero or more options.
func New(opts ...Option) Style {
	var s Style
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// Foreground sets the text color.
func Foreground(c Color) Option { return func(s *Style) { s.fg, s.hasFg = c, true } }

// Background sets the background color.
func Background(c Color) Option { return func(s *Style) { s.bg, s.hasBg = c, true } }

// Bold enables bold text.
func Bold() Option { return func(s *Style) { s.bold = true } }

// Underline enables underlined text.
func Underline() Option { return func(s *Style) { s.underline = true } }

// Render wraps text in this style's ANSI escape codes. It returns text
// unchanged when [Enabled] is false or the style is empty.
func (s Style) Render(text string) string {
	seq := s.sequence()
	if !Enabled || seq == "" {
		return text
	}
	return seq + text + reset
}

// Sprint renders the operands formatted with [fmt.Sprint].
func (s Style) Sprint(a ...any) string { return s.Render(fmt.Sprint(a...)) }

// Sprintf renders text formatted with [fmt.Sprintf].
func (s Style) Sprintf(format string, a ...any) string {
	return s.Render(fmt.Sprintf(format, a...))
}

// sequence returns the SGR escape sequence for the style, or "" if the style
// has no effect.
func (s Style) sequence() string {
	var params []string
	if s.bold {
		params = append(params, "1")
	}
	if s.underline {
		params = append(params, "4")
	}
	if s.hasFg {
		params = append(params, "38", "5", strconv.Itoa(int(s.fg)))
	}
	if s.hasBg {
		params = append(params, "48", "5", strconv.Itoa(int(s.bg)))
	}
	if len(params) == 0 {
		return ""
	}
	return csi + strings.Join(params, ";") + "m"
}

// --- JSON ---

type colorJSON struct {
	RGB   []uint8 `json:"rgb,omitempty"`
	Gray  *uint8  `json:"gray,omitempty"`
	Index *uint8  `json:"index,omitempty"`
	Name  string  `json:"name,omitempty"`
}

type styleJSON struct {
	Fg        *colorJSON `json:"fg,omitempty"`
	Bg        *colorJSON `json:"bg,omitempty"`
	Bold      bool       `json:"bold,omitempty"`
	Underline bool       `json:"underline,omitempty"`
}

var standardNames = [16]string{
	"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
	"bright-black", "bright-red", "bright-green", "bright-yellow",
	"bright-blue", "bright-magenta", "bright-cyan", "bright-white",
}

// toJSON picks the most readable representation that round-trips exactly.
func (c Color) toJSON() colorJSON {
	switch {
	case c >= 232:
		g := uint8(c - 232)
		return colorJSON{Gray: &g}
	case c >= 16:
		i := c - 16
		return colorJSON{RGB: []uint8{uint8(i / 36), uint8((i / 6) % 6), uint8(i % 6)}}
	default:
		return colorJSON{Name: standardNames[c]}
	}
}

func (cj colorJSON) toColor() (Color, error) {
	switch {
	case cj.Index != nil:
		return Color(*cj.Index), nil
	case cj.Gray != nil:
		return Gray(*cj.Gray), nil
	case len(cj.RGB) == 3:
		return RGB(cj.RGB[0], cj.RGB[1], cj.RGB[2]), nil
	case cj.Name != "":
		for i, n := range standardNames {
			if n == cj.Name {
				return Color(i), nil
			}
		}
		return 0, fmt.Errorf("colors: unknown color name %q", cj.Name)
	default:
		return 0, fmt.Errorf("colors: empty color spec; set one of rgb, gray, index, name")
	}
}

// MarshalJSON encodes the style as {"fg":…, "bg":…, "bold":…, "underline":…},
// omitting absent fields.
func (s Style) MarshalJSON() ([]byte, error) {
	sj := styleJSON{Bold: s.bold, Underline: s.underline}
	if s.hasFg {
		c := s.fg.toJSON()
		sj.Fg = &c
	}
	if s.hasBg {
		c := s.bg.toJSON()
		sj.Bg = &c
	}
	return json.Marshal(sj)
}

// UnmarshalJSON decodes the form produced by [Style.MarshalJSON]. A color may
// be given as {"rgb":[r,g,b]} (0–5), {"gray":n} (0–23), {"index":n} (0–255),
// or {"name":"red"}.
func (s *Style) UnmarshalJSON(b []byte) error {
	var sj styleJSON
	if err := json.Unmarshal(b, &sj); err != nil {
		return err
	}
	var out Style
	if sj.Fg != nil {
		c, err := sj.Fg.toColor()
		if err != nil {
			return err
		}
		out.fg, out.hasFg = c, true
	}
	if sj.Bg != nil {
		c, err := sj.Bg.toColor()
		if err != nil {
			return err
		}
		out.bg, out.hasBg = c, true
	}
	out.bold = sj.Bold
	out.underline = sj.Underline
	*s = out
	return nil
}
