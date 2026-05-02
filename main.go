package main

import (
	"strings"

	"github.com/anirban1809/tuix/tuix"
)

func nextLine() string {
	return "This is a line"
}

func App(props tuix.Props) tuix.Element {
	lines, setLines := tuix.UseState([]string{})

	if tuix.CurrentKey.Code == tuix.KeySpace {
		setLines(append(lines, nextLine()))
	}

	bodyStyle := tuix.NewStyle().Foreground(tuix.White).Bold(true).Background(tuix.Blue)
	accent := tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	// Full rounded border on all four sides, cyan-tinted.
	header := tuix.Box(
		tuix.Props{Direction: tuix.Row, Padding: [4]int{0, 2, 0, 2}},
		tuix.NewStyle().
			Background(tuix.Hex("#1e3a8a")).
			Border(tuix.Border{
				Top: true, Right: true, Bottom: true, Left: true,
				Chars: tuix.BorderRounded,
				Color: tuix.Cyan,
			}),
		tuix.Text("── spacebar appends a new line ──", accent),
	)

	longBlock := tuix.MultilineText(strings.Join(lines, "\n"), bodyStyle)

	paragraph := "The quick brown fox jumps over the lazy dog while a curious cat watches from the windowsill above."
	// Partial border: top + bottom rules only, no side edges.
	wrapped := tuix.Box(
		tuix.Props{Padding: [4]int{0, 1, 0, 1}},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Bottom: true,
			Chars: tuix.BorderSharp,
			Color: tuix.BrightYellow,
		}),
		tuix.WrappedText(paragraph, bodyStyle),
	)

	// Single-side accent: just a left rail.
	hint := tuix.Box(
		tuix.Props{Padding: [4]int{0, 0, 0, 1}},
		tuix.NewStyle().Border(tuix.Border{
			Left:  true,
			Chars: tuix.BorderThick,
			Color: tuix.BrightBlack,
		}),
		tuix.Text("Space to append a line · Esc to quit", dim),
	)

	// Outer Box: 1 row top/bottom, 2 cols left/right of gray padding around
	// the whole column. The padding is visible because the Box now paints
	// its own background.
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 4, 1, 4}},
		tuix.NewStyle().Background(tuix.Hex("#a1a1a1")),
		header,
		longBlock,
		wrapped,
		hint,
	)
}

func main() {
	app := tuix.NewApp(100, 10)
	app.Run(App, tuix.Props{})
}
