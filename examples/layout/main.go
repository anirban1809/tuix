// Example: layout
// Demonstrates: Direction, Gap, Padding, Sizing (Fit/Fixed/Grow), Align, Justify.
//
// You'll see a tiny dashboard mockup: a header that uses JustifySpaceBetween
// to push three labels to start/center/end of the row, a body that splits a
// fixed-width sidebar from a Grow(1) main pane, and a footer aligned right.
//
// Run: go run ./examples/layout
// See: ../../DOCS.md#layout

package main

import "github.com/anirban1809/tuix/tuix"

func App(props tuix.Props) tuix.Element {
	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	pane := tuix.NewStyle().Border(tuix.Border{
		Top: true, Right: true, Bottom: true, Left: true,
		Chars: tuix.BorderRounded, Color: tuix.BrightBlack,
	})

	header := tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
			Justify:   tuix.JustifySpaceBetween,
			Align:     tuix.AlignCenter,
		},
		tuix.NewStyle(),
		tuix.Text("◆ dashboard", title),
		tuix.Text("layout demo", dim),
		tuix.Text("v0.0.15", dim),
	)

	sidebar := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Fixed(20), Gap: 1},
		pane,
		tuix.Text("◌ Inbox", title),
		tuix.Text("  Archive", dim),
		tuix.Text("  Sent", dim),
		tuix.Text("  Drafts", dim),
	)

	main := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1), Gap: 1},
		pane,
		tuix.Text("Main pane (Grow(1))", title),
		tuix.WrappedText("This pane uses Grow(1) so it fills whatever width is left after the fixed-width sidebar takes its 20 columns. Resize your terminal to see it adapt.", dim),
	)

	body := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		sidebar,
		main,
	)

	footer := tuix.Box(
		tuix.Props{Direction: tuix.Row, Width: tuix.Grow(1), Justify: tuix.JustifyEnd},
		tuix.NewStyle(),
		tuix.Text("ctrl-c to quit", dim),
	)

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle(),
		header, body, footer,
	)
}

func main() {
	app := tuix.NewApp(90, 14)
	app.Run(App, tuix.Props{})
}
