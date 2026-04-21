package main

import (
	"fmt"
	"os"
	"time"

	"flowstacks.ai/tuix/tuix"
)

func Counter(props tuix.Props) []tuix.Element {
	label := fmt.Sprintf("Count: %d", props.Get("count"))
	return []tuix.Element{
		tuix.Box(tuix.Props{Direction: tuix.Row, Gap: 1},
			tuix.Text(label, tuix.Style{}),
			tuix.Text("[+]", tuix.Style{}.Foreground(tuix.Blue).Bold(true)),
		),
	}
}

func main() {
	screen := tuix.NewScreenWriter(80, 24, os.Stdout)
	r := tuix.NewRenderer(screen)

	screen.Start()
	defer screen.Stop()

	counter := 0

	for {
		// Frame 1: render count=0
		r.Render(tuix.Component(Counter, tuix.Props{Values: map[string]any{"count": counter}}))
		screen.Flush()
		time.Sleep(500 * time.Millisecond)
		counter++
		// Frame 2: re-render count=1 — reconciler must UPDATE not recreate
		r.Render(tuix.Component(Counter, tuix.Props{Values: map[string]any{"count": counter}}))
		screen.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}
