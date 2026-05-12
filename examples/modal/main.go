// Example: modal
// Demonstrates: opening and closing a Modal overlay.
// Enter to open · Esc to close.
//
// Run: go run ./examples/modal
// See: ../../DOCS.md#modal

package main

import (
	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	open, setOpen := tuix.UseState(false)

	if !open && tuix.CurrentKey.Code == tuix.KeyEnter {
		setOpen(true)
	}

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.BrightWhite)

	modal := components.Modal(
		"Are you sure?",
		open,
		36,
		func() { setOpen(false) },
		tuix.Text("This is a destructive action.", body),
		tuix.Text("Esc to cancel.", dim),
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ modal demo", title),
		tuix.Text("press Enter to open the modal", dim),
		modal,
	)
}

func main() {
	app := tuix.NewApp(60, 10)
	app.Run(App, tuix.Props{})
}
