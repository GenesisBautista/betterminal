# betterminal

A small, **zero-dependency** terminal styling toolkit for Go: 256-color text
styling, named semantic themes (loadable from JSON), and a console grid/table
renderer.

```
go get github.com/GenesisBautista/betterminal
```

## One import or three

The root package re-exports everything, so a **single import** covers colors,
themes, and tables:

```go
import "github.com/GenesisBautista/betterminal"

func main() {
	t := betterminal.NewGrid(betterminal.Padding(1, 1)).
		Header("Name", "Score").
		Row("Ada", 99).
		Row("Linus", 87)
	fmt.Print(t)
}
```

The facade is a thin alias layer (`NewGrid`, `RGB`, `Foreground`, `DefaultTheme`,
`RoleSuccess`, …). For finer-grained imports, the three subpackages are layered
low-to-high and can be used directly — the rest of this README uses them:

| Package  | Import                                          | What it does                                  |
| -------- | ----------------------------------------------- | --------------------------------------------- |
| `colors` | `github.com/GenesisBautista/betterminal/colors` | ANSI 256-color `Color` + `Style` (foreground/background, bold, underline) |
| `themes` | `github.com/GenesisBautista/betterminal/themes` | Named roles (Message, Success, Error, …) → styles, with JSON load/save |
| `grid`   | `github.com/GenesisBautista/betterminal/grid`   | Bordered, aligned, styled tables              |

Facade name mapping: `betterminal.NewStyle` = `colors.New`, `betterminal.NewGrid`
= `grid.New`, `betterminal.DefaultTheme` = `themes.Default`, roles are prefixed
(`betterminal.RoleError` = `themes.Error`). Everything else keeps its name.

Color output honors the [`NO_COLOR`](https://no-color.org) convention and can be
toggled with `colors.Enabled` (`betterminal.SetColorEnabled`).

---

## colors

A `Color` is an xterm-256 palette index. Build one three ways:

```go
colors.RGB(5, 0, 0) // 6×6×6 cube; each channel is a level 0–5
colors.Gray(8)      // 24-step grayscale ramp, 0 (dark) – 23 (light)
colors.Red          // named standard colors, indices 0–15
```

A `Style` is built with functional options and rendered onto a string:

```go
package main

import (
	"fmt"

	"github.com/GenesisBautista/betterminal/colors"
)

func main() {
	ok := colors.New(colors.Foreground(colors.RGB(0, 5, 0)), colors.Bold())
	warn := colors.New(
		colors.Foreground(colors.Black),
		colors.Background(colors.RGB(5, 5, 0)),
	)

	fmt.Println(ok.Render("build passed"))
	fmt.Println(warn.Sprintf("%d warnings", 3))
}
```

`Style` values are immutable and comparable. The zero `Style` renders text
unchanged.

### Turning color on/off

```go
colors.Enabled = colors.IsTerminal(os.Stdout) // only colorize a real terminal
```

`Enabled` defaults to `false` when `NO_COLOR` is set in the environment, `true`
otherwise. When it's `false`, `Render` returns the plain string.

---

## themes

A `Theme` maps semantic `Role`s to `colors.Style` values, so call sites refer to
intent ("this is an error") rather than a specific color.

```go
package main

import "github.com/GenesisBautista/betterminal/themes"

func main() {
	t := themes.Default()

	t.Println(themes.Success, "deploy complete")
	t.Println(themes.Error, "connection refused")
	fmt.Print(t.Sprint(themes.Info, "12 items\n"))
}
```

Roles: `Message`, `Success`, `Error`, `Warning`, `Info`, `System`, `Grid`
(iterate them with `themes.Roles`).

### Loading a theme from JSON

```go
t, err := themes.Load("theme.json")        // strict: returns the error
t, _ := themes.LoadOrDefault("theme.json") // falls back to themes.Default()
```

Each role is `{ "fg": <color>, "bg": <color>, "bold": bool, "underline": bool }`,
where a color is any one of:

```jsonc
{ "rgb": [5, 0, 0] }   // cube levels 0–5
{ "gray": 8 }          // ramp 0–23
{ "index": 196 }       // raw palette index 0–255
{ "name": "red" }      // a standard 0–15 name
```

Example `theme.json`:

```json
{
  "success": { "fg": { "rgb": [0, 5, 0] } },
  "error":   { "fg": { "rgb": [5, 0, 0] }, "bold": true },
  "info":    { "fg": { "rgb": [0, 0, 5] } }
}
```

Save a theme back out with `t.Save("theme.json")`.

---

## grid

Build a table by chaining `Header` and `Row`. Plain values are formatted with
`fmt.Sprint`; a `*Table` is a `fmt.Stringer`, so you can print it directly.

```go
package main

import (
	"fmt"

	"github.com/GenesisBautista/betterminal/grid"
)

func main() {
	t := grid.New(grid.Padding(1, 1)).
		Header("Name", "Score").
		Row("Ada", 99).
		Row("Linus", 87)

	fmt.Print(t)
}
```

```
┌───────┬───────┐
│ Name  │ Score │
├───────┼───────┤
│ Ada   │ 99    │
│ Linus │ 87    │
└───────┴───────┘
```

### Alignment, formatting, full grid

Alignment and color resolve with a **cell › column › table** precedence — a cell
override beats a column override beats the table default.

```go
t := grid.New(
	grid.Padding(1, 1),
	grid.Alignment(grid.AlignLeft), // table default
	grid.FullGrid(),                // border between every row
).
	Header("Symbol", "Price").
	Column(1,
		grid.ColumnAlign(grid.AlignRight),
		grid.ColumnFormat(func(v any) string { return fmt.Sprintf("$%.2f", v) }),
	).
	Row("AMZN", 556411.52).
	Row("TSLA", grid.C(560522.35, grid.CellAlign(grid.AlignLeft))) // cell overrides column
```

Per-column formatters apply to data cells only — header labels are left as-is.

### Styling a table with a theme

`grid` styles are just `colors.Style` values, so a theme drops straight in:

```go
theme := themes.Default()

t := grid.New(
	grid.HeaderStyle(theme.Style(themes.Info)),
	grid.BorderStyle(theme.Style(themes.Grid)),
	grid.AlternatingStyles(
		colors.New(colors.Foreground(colors.Gray(18))),
		colors.New(colors.Foreground(colors.Gray(12))),
	),
).
	Header("Symbol", "Side").
	Row("AAPL", grid.C("buy", grid.CellStyle(theme.Style(themes.Success))))
```

Use `grid.C(value, opts...)` to attach per-cell `CellAlign` / `CellStyle`.
Mismatched column counts don't panic — `t.Err()` reports the first problem while
the table still renders best-effort.

---

## Running the demo

`cmd/demo` exercises every package (run from the repo root so `theme.json`
resolves):

```
go run ./cmd/demo
```

## License

MIT — see [LICENSE](LICENSE).
