# betterminal

A terminal styling package for Go. Provides ANSI 256-color text styling, named semantic color themes (loadable from JSON), and — in progress — a console grid/table renderer modeled on the author's C# [ConsoleTableMaker](https://github.com/GenesisBautista/ConsoleTableMaker).

## Overview

- **Type**: Library (public Go module — consumed via `go get`)
- **Language**: Go 1.22.4
- **Framework**: none — standard library only
- **Database**: none
- **Testing**: table-driven tests with the stdlib `testing` package (`colors`, `themes`, `grid`)
- **Deploy target**: published as a Go module; no runtime/server deploy
- **Dependencies**: zero external — this is a hard constraint

## Commands

```bash
# Build everything
go build ./...

# Test
go test ./...

# Vet + format
go vet ./...
gofmt -l -w .

# Run the demo (exercises colors + themes + grid; run from repo root for theme.json)
go run ./cmd/demo
```

No code generation, no migrations.

## Architecture

A layered styling library, low-level to high-level:

1. **`colors/`** — the foundation. A `Color` is an xterm-256 palette index, built
   with `RGB(r,g,b)` (6×6×6 cube, levels 0–5), `Gray(level)` (0–23), or the named
   0–15 constants. A `Style` couples an optional foreground/background `Color` with
   attributes (bold, underline) and is constructed with **functional options**:
   `colors.New(colors.Foreground(c), colors.Bold())`. `Style.Render` wraps text in
   ANSI codes, gated by the package-level `Enabled` flag (NO_COLOR-aware).
2. **`themes/`** — a semantic layer on top of `colors`. `Theme` is a
   `map[Role]colors.Style` keyed by a string `Role` (Message, Success, Error,
   Warning, Info, System, Grid). `Default()` builds a full theme; `Load`/`Save`
   round-trip JSON; `Println`/`Sprint(role, …)` render in a role's style.
3. **`grid/`** — a console table renderer (idiomatic port of ConsoleTableMaker).
   `New(opts...)` then chain `Header`/`Row`; `*Table` is a `fmt.Stringer`. Cells
   carry optional align/style; alignment and color resolve **cell › column › table**.
   Depends only on `colors` (not `themes`), so a theme is passed in as plain
   `colors.Style` values.

On top sits the **root `betterminal` facade** (`betterminal.go`): a thin alias/wrapper
layer re-exporting the three subpackages so consumers can use a single import. It adds
no behavior.

Data flow: caller builds a `Style` (or looks one up by `Role` in a `Theme`, or builds
a `grid.Table`) → calls `Render`/`Sprintf`/`Println`/`String` → ANSI codes wrap the text
→ output to terminal.

## Directory Structure

```
betterminal.go One-import facade: package betterminal re-exports colors/themes/grid (aliases + wrappers).
colors/        Color (palette index) + Style (functional-options builder) + JSON. No deps.
themes/        map[Role]Style with Default/Load/Save + role-based print helpers; built on colors/.
grid/          Table renderer (Header/Row builder, borders, align/style hierarchy); built on colors/.
cmd/demo/      Runnable demo (package main) exercising every package; run with `go run ./cmd/demo`.
theme.json     Sample theme consumed by themes.Load (fg/bg as {rgb:[0-5]}, {gray:0-23}, {index}, {name}).
go.mod         Module github.com/GenesisBautista/betterminal (note capital G, B).
```

One package per concept; package name == directory name, lowercase, single word.

## Patterns & Conventions

- **Functional options** for public constructors: `New(opts ...Option)` with `WithX`/
  `Foreground`/`Background`/`Bold`/`Underline` option funcs. This is the house style —
  follow it for new packages (including the grid).
- **Zero external dependencies.** Solve with the stdlib. Adding a third-party package
  needs explicit sign-off.
- **Values over pointers** for small immutable types (`Color`, `Style` are value types
  and comparable — keep them so; it makes them map keys / `==`-testable).
- **Public API hygiene:** this ships as a public module. Exported identifiers get
  godoc comments; avoid breaking exported signatures without reason.
- **Keep the facade in sync:** when you add, rename, or remove an exported symbol in
  `colors`/`themes`/`grid`, mirror it in `betterminal.go`. Facade names: subpackage
  `New` becomes `NewStyle`/`NewGrid`, `themes.Default`→`DefaultTheme`, roles get a
  `Role` prefix (`themes.Error`→`RoleError`).

## Naming Conventions

- **Packages**: lowercase, single word, == directory name (`colors`, `themes`).
- **Types & exported funcs**: PascalCase (`Color`, `Style`, `Theme`, `Role`).
- **Constructors**: `New` for the options builder; `Default`/`Load` for whole-value
  factories.
- **Option funcs** are verbs/nouns describing the setting (`Foreground`, `Bold`).
- Raw ANSI string-building stays unexported behind methods (`sequence`, `csi`, `reset`).

## Error Handling

- Return errors up the call stack (standard Go), wrapped with `%w` and a package prefix
  (`fmt.Errorf("themes: read %s: %w", …)`).
- **Don't blur value + error:** `themes.Load` returns a clean `(Theme, error)`;
  `themes.LoadOrDefault` substitutes `Default()` on failure but still returns the error
  so callers can log it. Pick the one matching the call site's needs.
- No panics for expected failures (missing file, bad color spec).

## Testing

- Table-driven tests using the stdlib `testing` package. No assertion library
  (keeps the zero-dep rule).
- `colors`: `RGB`/`Gray` math + clamping, `Render` on/off via `Enabled`, exact SGR
  sequences, and JSON round-trip. Tests that toggle `colors.Enabled` must restore it
  with `defer`.
- `themes`: `Default` completeness, `Load`/`LoadOrDefault`, Save→Load round-trip.

## Environment Setup

- None — pure library. `go build ./...` works out of the box.
- Runtime requires a terminal that supports **ANSI 256-color**.
- Set `NO_COLOR` (any value) to disable color, per https://no-color.org; or set
  `colors.Enabled` directly (e.g. `colors.Enabled = colors.IsTerminal(os.Stdout)`).

## Gotchas

- **`colors.Enabled` is package-global mutable state.** Tests and demos flip it;
  always `defer` a restore so you don't leak the change into other tests.
- **Grayscale direction differs from the old API.** `Gray(0)` is now the *darkest*
  step (palette 232) and `Gray(23)` the lightest — the reverse of the original
  `BI_*` ordering. Watch this if porting old code/themes.
- **Case-sensitive import paths.** macOS's case-insensitive filesystem hides import-case
  mistakes; the module path is `github.com/GenesisBautista/betterminal` (capital G, B).
  Imports must match exactly or a Linux CI build will fail.
- **JSON color forms aren't symmetric on the wire.** Any of `{rgb}`, `{gray}`,
  `{index}`, `{name}` parse on input; on output `Style` emits the most readable exact
  form (gray→`gray`, cube→`rgb`, 0–15→`name`). It always round-trips to the same
  `Color`, but the JSON text may change shape.

## Off-Limits

- **Do not add external dependencies.** Stdlib only unless explicitly approved.
- Do not hand-edit `go.mod`/`go.sum` to pull in deps as a side effect.
- Do not edit `.git/` internals.

## Domain Language

**Color**: an xterm 256-color *palette index* (0–255), not a 24-bit RGB value.
Build cube colors with `RGB(r,g,b)` where each component is a **level 0–5**; build
ramp colors with `Gray(level)` where level is **0–23**.

**Style**: a foreground/background color plus attributes (bold, underline). The zero
`Style` renders text unchanged.

**Role**: a string-typed semantic slot in a `Theme` (`"message"`, `"success"`, …),
which doubles as the JSON key.

**Theme**: `map[Role]Style` — a set of named semantic roles, not a full UI theme.
A missing role yields the zero `Style` (safe, no styling).

**Grid**: a console table/grid renderer (the in-progress feature), modeled on the C#
ConsoleTableMaker.

## Roadmap / Project Direction

- ✅ **`colors` + `themes` refactored** to idiomatic Go: `Color`/`Style` with functional
  options, NO_COLOR-aware `Enabled`, `map[Role]Style` themes, clean JSON round-trip,
  and all roles (incl. System/Grid) covered. Tests added. Zero deps preserved.
- ✅ **Grid/table system built** (`grid/`), an idiomatic port of
  [ConsoleTableMaker](https://github.com/GenesisBautista/ConsoleTableMaker): box-drawing
  borders, padding, cell › column › table alignment + color hierarchy, alternating row
  styles, header separation, full-grid mode, per-column formatter funcs. Functional
  options, value types, tested, zero deps. README documents import + usage of all three
  packages.
- ⏭ Possible next steps: wide/CJK-aware width (currently rune count), per-cell column
  spanning, CSV/`[][]string` ingestion helpers, a `cmd/` example separate from root `main.go`.
