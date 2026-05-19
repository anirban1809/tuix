package main

import (
	"os"

	"github.com/anirban1809/tuix/tuix"
)

func App(props tuix.Props) tuix.Element {
	f, _ := os.ReadFile(
		"/Users/anirban/Documents/Code/zipcode-benchmarks/reports/20260518T195338Z/report.md",
	)

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle(),
		tuix.Markdown(string(f), tuix.NewStyle()),
	)
}

func main() {
	app := tuix.NewApp(80, 25)
	app.Run(App, tuix.Props{})
}
