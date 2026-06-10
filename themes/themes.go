// Package themes maps named semantic roles (Message, Success, Error, …) to
// reusable [colors.Style] values, and loads or saves them as JSON.
//
//	t := themes.Default()
//	t.Println(themes.Success, "done")
package themes

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GenesisBautista/betterminal/colors"
)

// Role names a semantic style slot in a [Theme]. Its string value is also the
// JSON key, so a theme round-trips to {"message": …, "success": …, …}.
type Role string

// The roles every [Default] theme defines.
const (
	Message Role = "message"
	Success Role = "success"
	Error   Role = "error"
	Warning Role = "warning"
	Info    Role = "info"
	System  Role = "system"
	Grid    Role = "grid"
)

// Roles lists the standard roles in display order.
var Roles = []Role{Message, Success, Error, Warning, Info, System, Grid}

// Theme maps roles to styles. A missing role yields the zero [colors.Style],
// which renders text unchanged, so lookups are always safe.
type Theme map[Role]colors.Style

// Default returns a theme with a sensible style for every role in [Roles].
func Default() Theme {
	return Theme{
		Message: colors.New(colors.Foreground(colors.Gray(12))),
		Success: colors.New(colors.Foreground(colors.RGB(0, 5, 0))),
		Error:   colors.New(colors.Foreground(colors.RGB(5, 0, 0)), colors.Bold()),
		Warning: colors.New(colors.Foreground(colors.RGB(5, 5, 0))),
		Info:    colors.New(colors.Foreground(colors.RGB(0, 0, 5))),
		System:  colors.New(colors.Foreground(colors.Gray(8))),
		Grid:    colors.New(colors.Foreground(colors.Gray(18))),
	}
}

// Style returns the style for role, or the zero style if role is unset.
func (t Theme) Style(role Role) colors.Style { return t[role] }

// Load reads and parses a theme from a JSON file. Unlike [LoadOrDefault] it
// reports any read or parse error instead of substituting defaults.
func Load(path string) (Theme, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("themes: read %s: %w", path, err)
	}
	var t Theme
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("themes: parse %s: %w", path, err)
	}
	return t, nil
}

// LoadOrDefault returns the theme at path, or [Default] if it cannot be read
// or parsed. The error is returned for logging but may be ignored.
func LoadOrDefault(path string) (Theme, error) {
	t, err := Load(path)
	if err != nil {
		return Default(), err
	}
	return t, nil
}

// Save writes the theme to path as indented JSON.
func (t Theme) Save(path string) error {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Sprint renders the operands in role's style using [fmt.Sprint] semantics.
func (t Theme) Sprint(role Role, a ...any) string {
	return t.Style(role).Render(fmt.Sprint(a...))
}

// Sprintf renders text in role's style using [fmt.Sprintf] semantics.
func (t Theme) Sprintf(role Role, format string, a ...any) string {
	return t.Style(role).Sprintf(format, a...)
}

// Fprintln renders the operands in role's style and writes them to w with a
// trailing newline.
func (t Theme) Fprintln(w io.Writer, role Role, a ...any) (int, error) {
	return fmt.Fprintln(w, t.Sprint(role, a...))
}

// Println renders the operands in role's style and writes them to standard
// output with a trailing newline.
func (t Theme) Println(role Role, a ...any) (int, error) {
	return t.Fprintln(os.Stdout, role, a...)
}
