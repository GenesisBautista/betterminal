package colors

import (
	"encoding/json"
	"testing"
)

func TestRGB(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		want    Color
	}{
		{"black corner", 0, 0, 0, 16},
		{"white corner", 5, 5, 5, 231},
		{"green", 0, 5, 0, 46},
		{"blue", 0, 0, 5, 21},
		{"mid", 2, 3, 4, 16 + 36*2 + 6*3 + 4},
		{"clamps high", 9, 9, 9, 231},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RGB(tt.r, tt.g, tt.b); got != tt.want {
				t.Errorf("RGB(%d,%d,%d) = %d, want %d", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestGray(t *testing.T) {
	tests := []struct {
		level uint8
		want  Color
	}{
		{0, 232}, {23, 255}, {10, 242}, {99, 255},
	}
	for _, tt := range tests {
		if got := Gray(tt.level); got != tt.want {
			t.Errorf("Gray(%d) = %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestRenderEnabled(t *testing.T) {
	defer restoreEnabled(Enabled)
	Enabled = true

	s := New(Foreground(Green), Bold())
	want := "\x1b[1;38;5;2mhi\x1b[0m"
	if got := s.Render("hi"); got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderDisabled(t *testing.T) {
	defer restoreEnabled(Enabled)
	Enabled = false

	s := New(Foreground(Green))
	if got := s.Render("hi"); got != "hi" {
		t.Errorf("Render() with Enabled=false = %q, want %q", got, "hi")
	}
}

func TestRenderZeroStyle(t *testing.T) {
	defer restoreEnabled(Enabled)
	Enabled = true

	if got := (Style{}).Render("plain"); got != "plain" {
		t.Errorf("zero Style.Render() = %q, want %q", got, "plain")
	}
}

func TestSequenceForegroundAndBackground(t *testing.T) {
	s := New(Foreground(RGB(0, 5, 0)), Background(Gray(3)), Underline())
	want := "\x1b[4;38;5;46;48;5;235m"
	if got := s.sequence(); got != want {
		t.Errorf("sequence() = %q, want %q", got, want)
	}
}

func TestStyleJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		style Style
	}{
		{"fg cube", New(Foreground(RGB(1, 2, 3)))},
		{"fg gray", New(Foreground(Gray(7)))},
		{"fg standard", New(Foreground(Red))},
		{"fg+bg+attrs", New(Foreground(Green), Background(Black), Bold(), Underline())},
		{"empty", Style{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.style)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			var got Style
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if got != tt.style {
				t.Errorf("round trip = %+v, want %+v (json %s)", got, tt.style, b)
			}
		})
	}
}

func TestUnmarshalColorForms(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Color
	}{
		{"rgb", `{"fg":{"rgb":[0,5,0]}}`, RGB(0, 5, 0)},
		{"gray", `{"fg":{"gray":7}}`, Gray(7)},
		{"index", `{"fg":{"index":200}}`, Color(200)},
		{"name", `{"fg":{"name":"magenta"}}`, Magenta},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Style
			if err := json.Unmarshal([]byte(tt.in), &s); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if !s.hasFg || s.fg != tt.want {
				t.Errorf("fg = %d (set=%v), want %d", s.fg, s.hasFg, tt.want)
			}
		})
	}
}

func TestUnmarshalBadColor(t *testing.T) {
	for _, in := range []string{`{"fg":{}}`, `{"fg":{"name":"chartreuse"}}`} {
		var s Style
		if err := json.Unmarshal([]byte(in), &s); err == nil {
			t.Errorf("Unmarshal(%s) = nil error, want error", in)
		}
	}
}

func restoreEnabled(v bool) { Enabled = v }
