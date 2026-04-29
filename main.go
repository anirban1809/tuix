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

	header := tuix.Text("── spacebar appends a new line ──", accent)

	longBlock := tuix.MultilineText(strings.Join(lines, "\n"), bodyStyle)

	hint := tuix.Text("Space to append a line · Esc to quit", dim)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle().Background(tuix.Hex("#a1a1a1")),
		header,
		longBlock,
		hint,
	)
}

func main() {
	app := tuix.NewApp(100, 10)
	app.Run(App, tuix.Props{})
}
