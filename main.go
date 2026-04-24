package main

import (
	"os"

	"flowstacks.ai/tuix/tuix"
	"flowstacks.ai/tuix/tuix/components"
)

// --- global state managed by the event loop ---
var (
	focusIdx = 0
	saved    = false
)

const numFocusable = 5

var (
	listItems     = []string{"Dashboard", "Settings", "Users", "Reports", "Logout"}
	selectOptions = []string{"Light", "Dark", "System", "High Contrast"}
)

func handleKey(key tuix.Key) {
	if key.Code == tuix.KeyTab {
		focusIdx = (focusIdx + 1) % numFocusable
		return
	}
	if key.Code == tuix.KeyShiftTab {
		focusIdx = (focusIdx - 1 + numFocusable) % numFocusable
		return
	}

	switch focusIdx {
	case 0: // Button — Enter toggles saved state
		if key.Code == tuix.KeyEnter {
			saved = !saved
		}

	case 1: // Input — TODO(human): handle text editing

	}
}

func App() tuix.Element {
	btnLabel := "Save Changes"
	if saved {
		btnLabel = "Saved ✓      "
	}

	dimStyle := tuix.Style{}.Foreground(tuix.BrightBlack)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},

		components.Panel("tuix — Interactive Demo", 52,
			components.Button(btnLabel, focusIdx == 0),
			tuix.Text(" ", tuix.Style{}),
			components.Input("Username: ", focusIdx == 1),
			components.Checkbox("Enable notifications", focusIdx == 2),
		),

		tuix.Box(
			tuix.Props{Direction: tuix.Row, Gap: 6},
			tuix.Box(
				tuix.Props{Direction: tuix.Column, Gap: 1},
				tuix.Text("Navigate", dimStyle),
				components.List(listItems, focusIdx == 3),
			),
			tuix.Box(
				tuix.Props{Direction: tuix.Column, Gap: 1},
				tuix.Text("Theme", dimStyle),
				components.SelectPicker(selectOptions, focusIdx == 4),
			),
		),

		tuix.Text("Tab · cycle focus    Shift+Tab · reverse    Esc · quit", dimStyle),
	)
}

func main() {
	tuix.StdOutScreen.Start()
	defer tuix.StdOutScreen.Stop()

	app := tuix.NewApp(100, 50)
	app.Run(App)

	quit := make(chan struct{})

	go func() {
		buf := make([]byte, 32)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				close(quit)
				return
			}
			key := tuix.ParseKey(buf[:n])
			if key.Code == tuix.KeyEscape || key.Code == tuix.KeyCtrlC {
				close(quit)
				return
			}
			tuix.Keys <- key
		}
	}()

	for {
		select {
		case <-quit:
			return
		case key := <-tuix.Keys:
			tuix.CurrentKey = key
			handleKey(key)
			app.Run(App)
		}
	}
}
