package themes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/GenesisBautista/betterminal/colors"
)

func TestDefaultHasEveryRole(t *testing.T) {
	d := Default()
	for _, r := range Roles {
		if _, ok := d[r]; !ok {
			t.Errorf("Default() missing role %q", r)
		}
	}
}

func TestStyleMissingRoleIsZero(t *testing.T) {
	if got := (Theme{}).Style(Grid); got != (colors.Style{}) {
		t.Errorf("Style of missing role = %+v, want zero", got)
	}
}

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme.json")
	data := `{
		"error": {"fg": {"rgb": [5,0,0]}, "bold": true},
		"system": {"fg": {"gray": 8}}
	}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if want := colors.New(colors.Foreground(colors.RGB(5, 0, 0)), colors.Bold()); got[Error] != want {
		t.Errorf("error role = %+v, want %+v", got[Error], want)
	}
	if want := colors.New(colors.Foreground(colors.Gray(8))); got[System] != want {
		t.Errorf("system role = %+v, want %+v", got[System], want)
	}
}

func TestLoadMissingFileErrors(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Error("Load(missing) = nil error, want error")
	}
}

func TestLoadOrDefaultFallsBack(t *testing.T) {
	got, err := LoadOrDefault(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Error("LoadOrDefault(missing) returned nil error, want the read error")
	}
	for _, r := range Roles {
		if _, ok := got[r]; !ok {
			t.Errorf("fallback theme missing role %q", r)
		}
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	if err := Default().Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Default()
	for _, r := range Roles {
		if got[r] != want[r] {
			t.Errorf("role %q = %+v, want %+v", r, got[r], want[r])
		}
	}
}

func TestSprintUsesStyle(t *testing.T) {
	defer func(v bool) { colors.Enabled = v }(colors.Enabled)
	colors.Enabled = false // assert plain text regardless of styling

	if got := Default().Sprint(Success, "ok"); got != "ok" {
		t.Errorf("Sprint = %q, want %q", got, "ok")
	}
}
