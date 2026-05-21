// Example: input
// Demonstrates: the Input component (typing, backspace, paste with Cmd+V),
// plus reading the live value back into a "preview" view.
//
// Run: go run ./examples/input
// See: ../../DOCS.md#input

package main

import (
	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	value, setValue := tuix.UseState("")

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.BrightWhite)

	field := tuix.Box(
		tuix.Props{Direction: tuix.Row, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded, Color: tuix.BrightYellow,
		}),
		// Input(label, focused, value, onChange).
		// focused=true here because there's only one field; in a real
		// form you'd track focus with another UseState.
		components.Input("name>", true, value, setValue),
	)

	preview := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{0, 1, 0, 1}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("you typed:", dim),
		tuix.WrappedText(value, body),
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ type something (paste works too!)", title),
		field,
		preview,
		tuix.Text("ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(80, 12)
	app.Run(App, tuix.Props{})
}
