// Example: styling
// Demonstrates: every color mode (ANSI16, ANSI256, Hex) and every border preset.
// Run:          go run ./examples/styling
// See:          ../../DOCS.md#styling

package main

import "github.com/anirban1809/tuix/tuix"

func swatch(label string, c tuix.Color) tuix.Element {
	return tuix.Text("  "+label+"  ",
		tuix.NewStyle().Background(c).Foreground(tuix.Black).Bold(true))
}

func borderedBox(label string, chars tuix.BorderChars, color tuix.Color) tuix.Element {
	return tuix.Box(
		tuix.Props{Padding: [4]int{0, 2, 0, 2}},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: chars,
			Color: color,
		}),
		tuix.Text(label, tuix.NewStyle().Foreground(color)),
	)
}

func App(props tuix.Props) tuix.Element {
	heading := tuix.NewStyle().Bold(true).Foreground(tuix.BrightWhite)

	ansi16 := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		swatch("red", tuix.Red),
		swatch("green", tuix.Green),
		swatch("yellow", tuix.Yellow),
		swatch("blue", tuix.Blue),
		swatch("magenta", tuix.Magenta),
		swatch("cyan", tuix.Cyan),
	)

	ansi256 := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		swatch("196", tuix.ANSI256(196)),
		swatch("208", tuix.ANSI256(208)),
		swatch("214", tuix.ANSI256(214)),
		swatch("220", tuix.ANSI256(220)),
		swatch("226", tuix.ANSI256(226)),
	)

	rgb := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 1},
		tuix.NewStyle(),
		swatch("#ff6b6b", tuix.Hex("#ff6b6b")),
		swatch("#ffd93d", tuix.Hex("#ffd93d")),
		swatch("#6bcf7f", tuix.Hex("#6bcf7f")),
		swatch("#4d9de0", tuix.Hex("#4d9de0")),
	)

	borders := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 2},
		tuix.NewStyle(),
		borderedBox("Sharp", tuix.BorderSharp, tuix.BrightWhite),
		borderedBox("Rounded", tuix.BorderRounded, tuix.Cyan),
		borderedBox("Double", tuix.BorderDouble, tuix.BrightYellow),
		borderedBox("Thick", tuix.BorderThick, tuix.BrightMagenta),
	)

	textStyles := tuix.Box(
		tuix.Props{Direction: tuix.Row, Gap: 2},
		tuix.NewStyle(),
		tuix.Text("bold", tuix.NewStyle().Bold(true)),
		tuix.Text("italic", tuix.NewStyle().Italic(true)),
		tuix.Text("underline", tuix.NewStyle().Underline(true)),
		tuix.Text("all three", tuix.NewStyle().Bold(true).Italic(true).Underline(true).Foreground(tuix.BrightCyan)),
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}},
		tuix.NewStyle(),
		tuix.Text("ANSI 16", heading), ansi16,
		tuix.Text("ANSI 256", heading), ansi256,
		tuix.Text("RGB / Hex", heading), rgb,
		tuix.Text("Border presets", heading), borders,
		tuix.Text("Text decorations", heading), textStyles,
	)
}

func main() {
	app := tuix.NewApp(80, 22)
	app.Run(App, tuix.Props{})
}
