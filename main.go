package main

import (
	"fmt"

	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	value, setValue := tuix.UseState("")
	lastPaste, setLastPaste := tuix.UseState("")
	pasteCount, setPasteCount := tuix.UseState(0)

	if tuix.CurrentKey.Code == tuix.KeyPaste {
		setLastPaste(tuix.CurrentKey.Paste)
		setPasteCount(pasteCount + 1)
	}

	title := tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.White)

	header := tuix.Box(
		tuix.Props{Direction: tuix.Row, Padding: [4]int{0, 1, 0, 1}},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded,
			Color: tuix.Cyan,
		}),
		tuix.Text("◆ bracketed paste demo — try Cmd+V", title),
	)

	field := tuix.Box(
		tuix.Props{Direction: tuix.Row, Padding: [4]int{0, 1, 0, 1}},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderSharp,
			Color: tuix.BrightYellow,
		}),
		components.Input("input>", "▌", true, value, setValue),
	)

	stats := tuix.Box(
		tuix.Props{Direction: tuix.Column},
		tuix.NewStyle(),
		tuix.Text(fmt.Sprintf("paste events: %d", pasteCount), dim),
		tuix.WrappedText("last paste (raw): "+lastPaste, body),
	)

	hint := tuix.Text("type to insert · backspace to delete · esc to quit", dim)

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle(),
		header,
		field,
		stats,
		hint,
	)
}

func main() {
	app := tuix.NewApp(100, 14)
	app.Run(App, tuix.Props{})
}
