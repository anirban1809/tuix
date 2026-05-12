// Example: hello
// Demonstrates: the smallest possible tuix program — one Box, one Text.
// Run:          go run ./examples/hello
// See:          ../../DOCS.md#quick-start

package main

import "github.com/anirban1809/tuix/tuix"

func App(props tuix.Props) tuix.Element {
	title := tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}},
		tuix.NewStyle(),
		tuix.Text("◆ hello, tuix", title),
		tuix.Text("(ctrl-c to quit)", dim),
	)
}

func main() {
	app := tuix.NewApp(60, 6)
	app.Run(App, tuix.Props{})
}
