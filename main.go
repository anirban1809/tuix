package main

import (
	"fmt"
	"strings"

	"github.com/anirban1809/tuix/tuix"
)

// buildLongText returns 60 numbered lines so the block is taller than any
// reasonable terminal viewport, ensuring the scroll-into-scrollback path
// in EnsureRoom is exercised.
func buildLongText() string {
	var b strings.Builder
	for i := 1; i <= 60; i++ {
		fmt.Fprintf(&b, "Line %02d ── scrollback verification line\n", i)
	}
	return b.String()
}

func App(props tuix.Props) tuix.Element {
	ticks, setTicks := tuix.UseState(0)

	// CurrentTick toggles every 500ms; count edges so we have a visibly
	// changing value at the bottom of the program.
	if tuix.CurrentTick {
		setTicks(ticks + 1)
	}

	bodyStyle := tuix.NewStyle().Foreground(tuix.White)
	accent := tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	header := tuix.Text("── multiline overflow demo ──", accent)

	longBlock := tuix.MultilineText(buildLongText(), bodyStyle)

	// Live counter — proves the visible portion of the program still
	// refreshes after older rows have scrolled into scrollback.
	live := tuix.Text(
		fmt.Sprintf("ticks: %d  (scroll up to see earlier lines)", ticks),
		accent,
	)

	hint := tuix.Text("Esc to quit", dim)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		header,
		longBlock,
		live,
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
