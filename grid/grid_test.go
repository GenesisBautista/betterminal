package grid

import (
	"strings"
	"testing"

	"github.com/GenesisBautista/betterminal/colors"
)

// withColorOff disables ANSI so rendered text can be compared exactly.
func withColorOff(t *testing.T) {
	t.Helper()
	prev := colors.Enabled
	colors.Enabled = false
	t.Cleanup(func() { colors.Enabled = prev })
}

func TestSimpleTable(t *testing.T) {
	withColorOff(t)

	got := New(Padding(1, 1)).
		Header("A", "B").
		Row("1", "22").
		String()

	want := strings.Join([]string{
		"┌───┬────┐",
		"│ A │ B  │",
		"├───┼────┤",
		"│ 1 │ 22 │",
		"└───┴────┘",
		"",
	}, "\n")

	if got != want {
		t.Errorf("table mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestNoHeadersNoSeparator(t *testing.T) {
	withColorOff(t)

	got := New().Row("x").Row("yy").String()
	want := strings.Join([]string{
		"┌──┐",
		"│x │",
		"│yy│",
		"└──┘",
		"",
	}, "\n")
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRightAlign(t *testing.T) {
	withColorOff(t)

	got := New(Alignment(AlignRight)).Row("a", "bbb").Row("cc", "d").String()
	// widths: col0=2, col1=3
	want := strings.Join([]string{
		"┌──┬───┐",
		"│ a│bbb│",
		"│cc│  d│",
		"└──┴───┘",
		"",
	}, "\n")
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestAlignmentHierarchyCellWins(t *testing.T) {
	withColorOff(t)

	// table=left, column 0 = right, but the cell forces center.
	got := New().
		Column(0, ColumnAlign(AlignRight)).
		Row(C("x", CellAlign(AlignCenter))).
		Row("yyy").
		String()
	want := strings.Join([]string{
		"┌───┐",
		"│ x │", // centered in width 3
		"│yyy│",
		"└───┘",
		"",
	}, "\n")
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestColumnFormatter(t *testing.T) {
	withColorOff(t)

	got := New().
		Column(0, ColumnFormat(func(v any) string { return "$" + fmtInt(v) })).
		Row(5).
		String()
	if !strings.Contains(got, "$5") {
		t.Errorf("formatter not applied:\n%s", got)
	}
}

func TestFullGridSeparators(t *testing.T) {
	withColorOff(t)

	got := New(FullGrid()).Row("a").Row("b").String()
	if strings.Count(got, "├") != 1 {
		t.Errorf("expected one inner separator, got:\n%s", got)
	}
}

func TestColumnCountMismatchSetsErr(t *testing.T) {
	tbl := New().Header("a", "b").Row("only-one")
	if tbl.Err() == nil {
		t.Error("Err() = nil, want a column-count mismatch error")
	}
	// still renders best-effort
	if tbl.String() == "" {
		t.Error("String() = empty, want best-effort render")
	}
}

func TestEmptyTable(t *testing.T) {
	if got := New().String(); got != "" {
		t.Errorf("empty table = %q, want empty", got)
	}
}

func TestStyledCellHasNoCodesWhenDisabled(t *testing.T) {
	withColorOff(t)

	got := New().Row(C("x", CellStyle(colors.New(colors.Foreground(colors.Red))))).String()
	if strings.Contains(got, "\x1b") {
		t.Errorf("escape codes leaked with color disabled:\n%q", got)
	}
}

func fmtInt(v any) string {
	if n, ok := v.(int); ok {
		return string(rune('0' + n))
	}
	return ""
}
