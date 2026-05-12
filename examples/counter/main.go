// Example: counter
// Demonstrates: UseState, keyboard events (Enter / Space / Backspace).
// Run:          go run ./examples/counter
// See:          ../../DOCS.md#usestate

package main

import (
	"fmt"

	"github.com/anirban1809/tuix/tuix"
)

func App(props tuix.Props) tuix.Element {
	count, setCount := tuix.UseState(0)

	switch tuix.CurrentKey.Code {
	case tuix.KeyEnter:
		setCount(count + 1)
	case tuix.KeySpace:
		setCount(count + 10)
	case tuix.KeyBackspace:
		setCount(0)
	}

	bigNumber := tuix.NewStyle().Bold(true).Foreground(tuix.BrightMagenta)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Align: tuix.AlignCenter},
		tuix.NewStyle(),
		tuix.Text(fmt.Sprintf("♥ %d", count), bigNumber),
		tuix.Text("enter +1 · space +10 · backspace reset · ctrl-c quit", dim),
	)
}

func main() {
	app := tuix.NewApp(70, 6)
	app.Run(App, tuix.Props{})
}
