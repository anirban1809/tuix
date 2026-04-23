package main

import (
	"fmt"
	"os"

	"flowstacks.ai/tuix/tuix"
)

// Counter renders a labeled count with independent state per instance.
func Counter(props tuix.Props) tuix.Element {
	label := props.Get("label").(string)
	step := props.Get("step").(int)

	count, setCount := tuix.UseState(5)
	setCount(count + step)

	return tuix.Text(
		fmt.Sprintf("%s: %d", label, count),
		tuix.Style{}.Foreground(tuix.Cyan),
	)
}

// Toggle tracks a boolean on/off state independently from the counters.
func Toggle(props tuix.Props) tuix.Element {
	on, setOn := tuix.UseState(false)

	setOn(!on)

	label := "OFF"
	style := tuix.Style{}.Foreground(tuix.Red)
	if on {
		label = "ON"
		style = tuix.Style{}.Foreground(tuix.Green)
	}

	return tuix.Text(fmt.Sprintf("Toggle: %s", label), style)
}

func App() tuix.Element {
	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		Counter(tuix.Props{Values: map[string]any{"label": "A", "step": 1}}),
		Counter(tuix.Props{Values: map[string]any{"label": "B", "step": 3}}),
		Toggle(tuix.Props{}),
	)
}

func main() {
	tuix.StdOutScreen.Start()
	defer tuix.StdOutScreen.Stop()

	app := tuix.NewApp(60, 10)
	app.Run(App)

	buf := make([]byte, 32)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || (n == 1 && buf[0] == 'q') {
			break
		}
		app.Run(App)
	}
}
