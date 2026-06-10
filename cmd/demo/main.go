// Command demo exercises every betterminal package. Run it from the repo root
// so the relative theme.json path resolves:
//
//	go run ./cmd/demo
package main

import (
	"fmt"
	"os"

	"github.com/GenesisBautista/betterminal"
	"github.com/GenesisBautista/betterminal/colors"
	"github.com/GenesisBautista/betterminal/grid"
	"github.com/GenesisBautista/betterminal/themes"
)

func main() {
	colors.Enabled = true // force color for the demo, even when piped

	colorDemo()
	grayscaleDemo()
	themeDemo()
	themeJSONDemo()
	gridDemo()
	facadeDemo()
}

func colorDemo() {
	red := colors.New(colors.Foreground(colors.RGB(5, 0, 0)))
	fmt.Println(red.Render("Red text"))

	redOnCyan := colors.New(
		colors.Foreground(colors.RGB(5, 0, 0)),
		colors.Background(colors.RGB(0, 5, 5)),
	)
	fmt.Println(redOnCyan.Render("Red text on a cyan background"))
}

func grayscaleDemo() {
	fmt.Println(colors.New(colors.Foreground(colors.Gray(5))).Render("Grayscale text"))

	for level := uint8(0); level <= 23; level++ {
		style := colors.New(colors.Background(colors.Gray(level)))
		fmt.Println(style.Sprintf("grayscale background level %d", level))
	}
}

func themeDemo() {
	t := themes.Default()
	t.Println(themes.Message, "This is a message.")
	t.Println(themes.Success, "This is a success message.")
	t.Println(themes.Error, "This is an error message.")
	t.Println(themes.Warning, "This is a warning message.")
	t.Println(themes.Info, "This is an info message.")
	t.Println(themes.System, "This is a system message.")
	t.Println(themes.Grid, "This is grid-colored text.")
}

func themeJSONDemo() {
	t, err := themes.Load("theme.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load theme:", err)
		return
	}
	for _, role := range themes.Roles {
		t.Println(role, "JSON theme role:", role)
	}
}

func gridDemo() {
	theme := themes.Default()

	table := grid.New(
		grid.Padding(1, 1),
		grid.HeaderStyle(theme.Style(themes.Info)),
		grid.BorderStyle(theme.Style(themes.Grid)),
		grid.AlternatingStyles(
			colors.New(colors.Foreground(colors.Gray(18))),
			colors.New(colors.Foreground(colors.Gray(12))),
		),
	).
		Header("Symbol", "Side", "Price", "Date").
		Column(2, grid.ColumnAlign(grid.AlignRight),
			grid.ColumnFormat(func(v any) string { return fmt.Sprintf("$%.2f", v) })).
		Row("AMZN", "buy", 556411.52, "2021-03-16").
		Row("TSLA", "sell", 560522.35, "2021-10-23").
		Row("AAPL", grid.C("buy", grid.CellStyle(theme.Style(themes.Success))), 755760.92, "2021-09-23")

	fmt.Print(table)
	if err := table.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "grid:", err)
	}
}

// facadeDemo builds the same kind of table using only the root betterminal
// package — a single import covers colors, themes, and grid.
func facadeDemo() {
	t := betterminal.NewGrid(betterminal.Padding(1, 1)).
		Header("Package", "Note").
		Row("colors", betterminal.C("low-level", betterminal.CellStyle(
			betterminal.NewStyle(betterminal.Foreground(betterminal.Green))))).
		Row("grid", "built on colors")
	fmt.Print(t)
}
