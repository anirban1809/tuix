// Example: conditional
// Demonstrates: the If helper for swapping between two pre-built elements.
// Space toggles between a "logged in" and "logged out" view.
//
// If(condition, a, b) is a regular function call, so both branches are
// evaluated before If runs. Pass pre-constructed elements — don't try to
// hide expensive work behind one branch.
//
// Run: go run ./examples/conditional
// See: ../../DOCS.md#conditional-rendering-with-if

package main

import "github.com/anirban1809/tuix/tuix"

func App(props tuix.Props) tuix.Element {
	loggedIn, setLoggedIn := tuix.UseState(false)
	if tuix.CurrentKey.Code == tuix.KeySpace {
		setLoggedIn(!loggedIn)
	}

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	greeting := tuix.Text("Welcome back, Avery 👋", title)
	prompt := tuix.Text("Please sign in.", tuix.NewStyle().Foreground(tuix.BrightYellow))

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Align: tuix.AlignCenter},
		tuix.NewStyle(),
		tuix.If(loggedIn, greeting, prompt),
		tuix.Text("space to toggle · ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(60, 6)
	app.Run(App, tuix.Props{})
}
