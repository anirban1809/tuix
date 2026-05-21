package main

import (
	"github.com/anirban1809/tuix/tuix"
)

func App(props tuix.Props) tuix.Element {
	const sample = `# Markdown example

Normal list:

- parsed bullet
1. parsed ordered item

Indented output:

    - not a list
    1. not ordered
`

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle(),
		tuix.Markdown(sample, tuix.NewStyle()),
	)
}

func main() {
	app := tuix.NewApp(80, 25)
	app.Run(App, tuix.Props{})
}
