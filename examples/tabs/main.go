// Example: tabs
// Demonstrates: the Tabs component switching content panels.
// Left/Right to switch tabs.
//
// Run: go run ./examples/tabs
// See: ../../DOCS.md#tabs

package main

import (
	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	tabs := []string{"Overview", "Activity", "Settings"}
	active, setActive := tuix.UseState(0)

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	body := tuix.NewStyle().Foreground(tuix.BrightWhite)

	var panel tuix.Element
	switch active {
	case 0:
		panel = tuix.WrappedText("Welcome back! You have 3 unread notifications and 1 PR awaiting review.", body)
	case 1:
		panel = tuix.WrappedText("· merged feat/context (2h ago)\n· opened docs/comprehensive-rewrite (5m ago)\n· reviewed #142 (yesterday)", body)
	case 2:
		panel = tuix.WrappedText("Theme: dark · Notifications: on · Telemetry: off", body)
	}

	bar := components.Tabs(tabs, true, setActive)
	pane := tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{1, 2, 1, 2}, Width: tuix.Grow(1)},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded, Color: tuix.BrightBlack,
		}),
		panel,
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text("◆ tabbed view", title),
		bar,
		pane,
		tuix.Text("←/→ to switch tabs · ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(80, 14)
	app.Run(App, tuix.Props{})
}
