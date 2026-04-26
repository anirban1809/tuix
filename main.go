package main

import (
	"strings"

	"github.com/anirban1809/tuix/tuix"
)

// nextLine returns the content of the line that should be appended when
// the user presses spacebar. `count` is the number of lines already in
// the block (so the first call receives 0).
func nextLine() string {
	return "This is a line"
}

func App(props tuix.Props) tuix.Element {
	lines, setLines := tuix.UseState([]string{})

	if tuix.CurrentKey.Code == tuix.KeySpace {
		setLines(append(lines, nextLine()))
	}

	bodyStyle := tuix.NewStyle().Foreground(tuix.White)
	accent := tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	header := tuix.Text("── spacebar appends a new line ──", accent)

	longBlock := tuix.MultilineText(strings.Join(lines, "\n"), bodyStyle)

	hint := tuix.Text("Space to append a line · Esc to quit", dim)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		header,
		longBlock,
		hint,
	)
}

func main() {
	// Width 60 keeps every painted row under the typical 80-col terminal
	// so EnsureRoom's inline writes don't wrap. Initial height 10 is
	// deliberately smaller than contentH so the scroll path runs on
	// frame 1.
	app := tuix.NewApp(60, 10)
	app.Run(App, tuix.Props{})
}
