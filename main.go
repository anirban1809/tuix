package main

import (
	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

// Theme is the value carried by ThemeContext. Components anywhere in the
// tree call tuix.UseContext(ThemeContext) to read the current theme
// without having to receive it through every intermediate ancestor.
type Theme struct {
	Name   string
	Fg     tuix.Color
	Bg     tuix.Color
	Accent tuix.Color
}

var lightTheme = Theme{
	Name:   "light",
	Fg:     tuix.Black,
	Bg:     tuix.BrightWhite,
	Accent: tuix.Blue,
}

var darkTheme = Theme{
	Name:   "dark",
	Fg:     tuix.BrightWhite,
	Bg:     tuix.Black,
	Accent: tuix.BrightCyan,
}

// ThemeContext is created at package scope so every component refers to
// the same Context identity. The default value (lightTheme) is what
// UseContext returns when no enclosing Provide is active.
var ThemeContext = tuix.CreateContext(lightTheme)

func Header() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle().Foreground(t.Accent).Bold(true),
		tuix.Text("◆ tuix context demo", tuix.Style{}),
	)
}

// Card sits two levels deep, demonstrating that context flows through
// any number of intermediate Box layers without prop-drilling.
func Card(title, body string) tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle().Foreground(t.Fg).Background(t.Bg).Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded,
			Color: t.Accent,
		}),
		tuix.Text(title, tuix.NewStyle().Foreground(t.Accent).Bold(true)),
		tuix.Text(body, tuix.NewStyle().Foreground(t.Fg).Background(t.Bg)),
	)
}

func Footer() tuix.Element {
	t := tuix.UseContext(ThemeContext)
	return tuix.Text(
		"theme: "+t.Name+" · t to toggle · q to quit",
		tuix.NewStyle().Foreground(t.Accent),
	)
}

func InputDemo() tuix.Element {
	value, setValue := tuix.UseState("edit me")

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Width: tuix.Grow(1)},
		tuix.NewStyle(),
		tuix.Text(
			"Input demo: use ←/→ to move cursor; type/paste inserts at cursor",
			tuix.NewStyle().Foreground(tuix.BrightBlack),
		),
		components.Input(">", true, value, setValue),
	)
}

func App(props tuix.Props) tuix.Element {
	dark, setDark := tuix.UseState(false)
	if tuix.CurrentKey.Rune == 't' {
		setDark(!dark)
	}
	if tuix.CurrentKey.Rune == 'q' {
		tuix.Exit()
	}

	current := lightTheme
	if dark {
		current = darkTheme
	}

	// Provide takes a render thunk rather than pre-built children. The
	// thunk runs *while* `current` is on the context stack, so every
	// Header/Card/Footer call evaluated inside it sees the active
	// theme via UseContext. Children built outside the thunk would
	// execute before the push and miss the value.
	return ThemeContext.Provide(current, func() tuix.Element {
		return tuix.Box(
			tuix.Props{
				Direction: tuix.Column,
				Gap:       1,
				Padding:   [4]int{1, 2, 1, 2},
				Width:     tuix.Grow(1),
			},
			tuix.NewStyle(),
			Header(),
			Card(
				"State of the world",
				"Both cards below consume the same ThemeContext.",
			),
			Card(
				"Sibling card",
				"Toggle the theme; both update without taking a theme prop.",
			),
			InputDemo(),
			Footer(),
		)
	})
}

func main() {
	app := tuix.NewApp(100, 18)
	app.Run(App, tuix.Props{})
}
