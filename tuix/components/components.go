package components

import (
	"strings"

	"flowstacks.ai/tuix/tuix"
)

// Panel renders a bordered container with a title in the top edge.
func Panel(title string, width int, children ...tuix.Element) tuix.Element {

	inner := width - 2
	titlePart := "─ " + title + " "
	remaining := max(inner-len([]rune(titlePart)), 0)

	top := "┌" + titlePart + strings.Repeat("─", remaining) + "┐"
	bottom := "└" + strings.Repeat("─", inner) + "┘"

	rows := []tuix.Element{tuix.Text(top, tuix.Style{})}
	for _, child := range children {
		rows = append(rows, tuix.Box(
			tuix.Props{Direction: tuix.Row},
			tuix.Text("│ ", tuix.Style{}),
			child,
		))
	}
	rows = append(rows, tuix.Text(bottom, tuix.Style{}))

	return tuix.Box(tuix.Props{Direction: tuix.Column}, rows...)
}

// Badge renders a short colored label with padding.
func Badge(label string, fg tuix.Color, bg tuix.Color) tuix.Element {
	return tuix.Text(
		" "+label+" ",
		tuix.Style{}.Foreground(fg).Background(bg).Bold(true),
	)
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner renders an animated braille spinner. Advances one frame per render.
func Spinner(label string) tuix.Element {
	frame, setFrame := tuix.UseState(0)
	setFrame((frame + 1) % len(spinnerFrames))

	return tuix.Text(
		spinnerFrames[frame]+" "+label,
		tuix.Style{}.Foreground(tuix.Cyan),
	)
}

// ProgressBar renders a filled bar. value must be between 0.0 and 1.0.
func ProgressBar(value float64, width int, color tuix.Color) tuix.Element {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}

	inner := width - 2
	filled := int(float64(inner) * value)
	empty := inner - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return tuix.Text(bar, tuix.Style{}.Foreground(color))
}

type AlertKind int

const (
	AlertInfo AlertKind = iota
	AlertSuccess
	AlertWarning
	AlertError
)

// Alert renders a status message with an icon prefix and matching color.
func Alert(kind AlertKind, message string) tuix.Element {
	var icon string
	var style tuix.Style

	switch kind {
	case AlertInfo:
		icon = "ℹ"
		style = tuix.Style{}.Foreground(tuix.Cyan)
	case AlertSuccess:
		icon = "✓"
		style = tuix.Style{}.Foreground(tuix.Green)
	case AlertWarning:
		icon = "⚠"
		style = tuix.Style{}.Foreground(tuix.Yellow)
	case AlertError:
		icon = "✗"
		style = tuix.Style{}.Foreground(tuix.Red)
	}

	return tuix.Text(icon+" "+message, style)
}
