package tuix_test

import (
	"strings"
	"testing"

	tuix "flowstacks.ai/tuix/tuix"
)

func TestStyleZeroValueIsValid(t *testing.T) {
	var s tuix.Style
	_ = s.ANSIPrefix()
}

func TestStyleIsValueType(t *testing.T) {
	a := tuix.Style{}.Bold(true).Foreground(tuix.Red)
	b := a.Bold(false)

	if !a.IsBold() {
		t.Error("original style was mutated")
	}
	if b.IsBold() {
		t.Error("derived style should not be bold")
	}
}

func TestStyleEquality(t *testing.T) {
	a := tuix.Style{}.Bold(true).Foreground(tuix.Red)
	b := tuix.Style{}.Bold(true).Foreground(tuix.Red)
	c := tuix.Style{}.Bold(true).Foreground(tuix.Blue)

	if a != b {
		t.Error("identical styles should be equal")
	}
	if a == c {
		t.Error("different styles should not be equal")
	}
}

func TestStyleANSIOutput(t *testing.T) {
	cases := []struct {
		name     string
		style    tuix.Style
		contains string
	}{
		{"bold", tuix.Style{}.Bold(true), "\033[1m"},
		{"italic", tuix.Style{}.Italic(true), "\033[3m"},
		{"underline", tuix.Style{}.Underline(true), "\033[4m"},
		{"fg red ansi16", tuix.Style{}.Foreground(tuix.Red), "\033[31m"},
		{"fg ansi256", tuix.Style{}.Foreground(tuix.ANSI256(200)), "\033[38;5;200m"},
		{"fg truecolor", tuix.Style{}.Foreground(tuix.Hex("#FF6B6B")), "\033[38;2;255;107;107m"},
		{"bg ansi16", tuix.Style{}.Background(tuix.Blue), "\033[44m"},
		{"bg ansi256", tuix.Style{}.Background(tuix.ANSI256(236)), "\033[48;5;236m"},
		{"bg truecolor", tuix.Style{}.Background(tuix.Hex("#1E1E2E")), "\033[48;2;30;30;46m"},
		{"reset", tuix.Style{}, "\033[0m"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.ANSIPrefix()
			if !strings.Contains(got, tt.contains) {
				t.Errorf("ANSIPrefix() = %q, want it to contain %q", got, tt.contains)
			}
		})
	}
}

func TestHexColorParsing(t *testing.T) {
	cases := []struct {
		hex     string
		r, g, b uint8
	}{
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#1E1E2E", 30, 30, 46},
		{"#ffffff", 255, 255, 255},
	}

	for _, tt := range cases {
		t.Run(tt.hex, func(t *testing.T) {
			c := tuix.Hex(tt.hex)
			r, g, b := c.RGB()
			if r != tt.r || g != tt.g || b != tt.b {
				t.Errorf("Hex(%q) = (%d,%d,%d), want (%d,%d,%d)", tt.hex, r, g, b, tt.r, tt.g, tt.b)
			}
		})
	}
}
