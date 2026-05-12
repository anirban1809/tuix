// Example: table
// Demonstrates: the Table component with onChange to read selected row.
// Up/Down navigates rows. The selected row's name is shown below the table.
//
// Run: go run ./examples/table
// See: ../../DOCS.md#table

package main

import (
	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

func App(props tuix.Props) tuix.Element {
	headers := []string{"Rank", "Player", "Score", "Streak"}
	rows := [][]string{
		{"1", "Riley", "9420", "12"},
		{"2", "Jules", "8810", "7"},
		{"3", "Theo", "8730", "9"},
		{"4", "Sam", "8112", "3"},
		{"5", "Avery", "7790", "5"},
	}

	selected, setSelected := tuix.UseState(0)

	title := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)
	hi := tuix.NewStyle().Foreground(tuix.BrightYellow)

	leaderboard := components.Table(headers, rows, true, setSelected)
	caption := "you are looking at: " + rows[selected][1]

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}},
		tuix.NewStyle(),
		tuix.Text("◆ leaderboard", title),
		leaderboard,
		tuix.Text(caption, hi),
		tuix.Text("↑/↓ to navigate · ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(70, 16)
	app.Run(App, tuix.Props{})
}
