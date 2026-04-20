package main

import (
	"os"

	"flowstacks.ai/tuix/tuix"
	"golang.org/x/term"
)

func main() {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	screen := tuix.NewScreenWriter(width, height, os.Stdout)
	screen.Start()
	defer screen.Stop()

	// Build a layout tree:
	// Root column
	//   ├── Fixed-height header row (height=3)
	//   └── Grow(1) body row
	//         ├── Fixed-width sidebar (width=20)
	//         └── Grow(1) main content area

	root := tuix.NewLayout().
		WithDirection(tuix.Column).
		WithSize(tuix.Fixed(width), tuix.Fixed(height)).
		WithChildren(
			tuix.NewLayout().
				WithDirection(tuix.Row).
				WithSize(tuix.Grow(1), tuix.Fixed(3)), // header

			tuix.NewLayout().
				WithDirection(tuix.Row).
				WithSize(tuix.Grow(1), tuix.Grow(1)). // body
				WithChildren(
					tuix.NewLayout().
						WithSize(tuix.Fixed(20), tuix.Grow(1)), // sidebar
					tuix.NewLayout().
						WithSize(tuix.Grow(1), tuix.Grow(1)), // main
				),
		)

	rects := tuix.ComputeLayout(root, tuix.Rect{X: 0, Y: 0, Width: width, Height: height})

	// Draw each rect as a colored border so you can see the boxes
	colors := []tuix.Color{tuix.Red, tuix.Green, tuix.Blue, tuix.Yellow}
	for i, r := range rects {
		style := tuix.Style{}.Foreground(colors[i%len(colors)])
		drawBorder(screen, r, style)
	}

	screen.Flush()

	// Wait for keypress then exit
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
}

func drawBorder(s *tuix.Screen, r tuix.Rect, style tuix.Style) {
	for x := r.X; x < r.X+r.Width; x++ {
		s.SetCell(x, r.Y, '─', style)
		s.SetCell(x, r.Y+r.Height-1, '─', style)
	}
	for y := r.Y; y < r.Y+r.Height; y++ {
		s.SetCell(r.X, y, '│', style)
		s.SetCell(r.X+r.Width-1, y, '│', style)
	}
	s.SetCell(r.X, r.Y, '┌', style)
	s.SetCell(r.X+r.Width-1, r.Y, '┐', style)
	s.SetCell(r.X, r.Y+r.Height-1, '└', style)
	s.SetCell(r.X+r.Width-1, r.Y+r.Height-1, '┘', style)
}
