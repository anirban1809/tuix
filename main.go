package main

import (
	"fmt"

	"flowstacks.ai/tuix/tuix"
)

func update[T any](p *T, v T) {
	*p = v
}

func Banner(props tuix.Props) []tuix.Element {
	title := "Hello world"
	var result []tuix.Element
	result = append(result, tuix.Text(title, tuix.Style{}.Foreground(tuix.Blue)), Counter(props)[0])
	return result
}

func Counter(props tuix.Props) []tuix.Element {
	count := 0
	update(&count, 1)

	return []tuix.Element{
		tuix.Box(tuix.Props{Direction: tuix.Row, Padding: [4]int{1, 1, 1, 1}},
			tuix.Text(fmt.Sprintf("Count: %d", count), tuix.Style{}),
			tuix.Text("[+]", tuix.Style{}.Foreground(tuix.Blue).Bold(true)),
		),
	}
}

func main() {
	tuix.StdOutScreen.Start()
	defer tuix.StdOutScreen.Stop()

	app := tuix.NewApp(80, 24)
	app.Run(Banner)
}
