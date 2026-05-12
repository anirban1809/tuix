// Example: list
// Demonstrates: the List component. Up/Down to move the highlight.
// Note: components.List manages its own selection state internally and does
// not currently surface the selected index back to the caller. Use it for
// pure visual selection; for programmatic access track your own cursor.
//
// Run: go run ./examples/list
// See: ../../DOCS.md#list

package main

import "github.com/anirban1809/tuix/tuix"
import "github.com/anirban1809/tuix/tuix/components"

func App(props tuix.Props) tuix.Element {
	playlist := []string{
		"♪ Strobe — deadmau5",
		"♪ Teardrop — Massive Attack",
		"♪ Midnight City — M83",
		"♪ Intro — The xx",
		"♪ Sunset Lover — Petit Biscuit",
		"♪ Resonance — Home",
	}

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightMagenta)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}},
		tuix.NewStyle(),
		tuix.Text("◆ now playing", title),
		components.List(playlist, true),
		tuix.Text("↑/↓ to browse · ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(70, 14)
	app.Run(App, tuix.Props{})
}
