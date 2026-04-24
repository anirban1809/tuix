package main

import (
	"fmt"

	"flowstacks.ai/tuix/tuix"
	"flowstacks.ai/tuix/tuix/components"
)

func App() tuix.Element {
	count, setCount := tuix.UseState(0)
	milestone, setMilestone := tuix.UseState("Press ↑ to start counting...")

	if tuix.CurrentKey.Code == tuix.KeyUp {
		setCount(count + 1)
	}
	if tuix.CurrentKey.Code == tuix.KeyDown && count > 0 {
		setCount(count - 1)
	}

	// Fires only when count changes; updates milestone message as a side effect.
	tuix.UseEffect(func() func() {
		if count > 0 && count%5 == 0 {
			setMilestone(fmt.Sprintf("Milestone reached: %d! Keep going...", count))
		}
		return nil
	}, []any{count})

	// Progress toward next multiple of 5
	progress := float64(count%5) / 5.0

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		tuix.Text(fmt.Sprintf("Count: %d  (↑/↓ to change, Esc to quit)", count), tuix.NewStyle()),
		components.ProgressBar(progress, 30, tuix.Green),
		components.Alert(components.AlertSuccess, milestone),
	)
}

func main() {
	app := tuix.NewApp(100, 10)
	app.Run(App)
}
