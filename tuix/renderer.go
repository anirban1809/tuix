package tuix

import (
	"strings"
)

// rawLines returns the text split on '\n' only, with no word wrapping.
func rawLines(text string) []string {
	return strings.Split(text, "\n")
}

// wrappedLines splits the text on '\n' and then word-wraps each segment to
// fit within maxWidth columns.
func wrappedLines(text string, maxWidth int) []string {
	var out []string
	for _, seg := range rawLines(text) {
		out = append(out, wrapText(seg, maxWidth)...)
	}
	return out
}

// wrapText breaks a single line of text (no '\n') into one or more lines so
// that each line's total cell width (per RuneWidth) is <= maxWidth. Breaks
// should prefer whitespace boundaries.
func wrapText(text string, maxWidth int) []string {
	lines := []string{}

	var line strings.Builder
	for i := 0; i < len(text); i++ {
		if i > 0 && i%maxWidth == 0 {
			lines = append(lines, line.String())
			line.Reset()
		}
		line.WriteByte(text[i])
	}

	lines = append(lines, line.String())
	return lines
}

type ComponentRenderer struct {
	screen *Screen
	dirty  chan struct{}
}

func NewRenderer(screen *Screen) *ComponentRenderer {
	return &ComponentRenderer{screen: screen, dirty: make(chan struct{})}
}

func (r *ComponentRenderer) Render(next Element) {
	layoutRoot := buildLayoutTree(next)

	_, contentH := IntrinsicSize(layoutRoot)
	screenH := r.screen.Height()
	if contentH > screenH {
		r.screen.Resize(r.screen.Width(), contentH)
	}

	available := Rect{X: 0, Y: 0, Width: r.screen.Width(), Height: r.screen.Height()}
	rects := ComputeLayout(layoutRoot, available)

	r.screen.Clear()
	paint(next, rects, 0, r.screen, Style{})

	// After paint: if contentH overflows the terminal viewport, write the
	// rows inline so the terminal scrolls older content into scrollback.
	// Must run after paint because it reads from the cell grid.
	r.screen.EnsureRoom(contentH)
}

func buildLayoutTree(element Element) *LayoutNode {
	p := element.Layout
	b := element.Style.border

	padTop, padRight, padBottom, padLeft := p.PaddingTop, p.PaddingRight, p.PaddingBottom, p.PaddingLeft
	if b.Top {
		padTop++
	}
	if b.Right {
		padRight++
	}
	if b.Bottom {
		padBottom++
	}
	if b.Left {
		padLeft++
	}

	l := &LayoutNode{
		Direction:     p.Direction,
		WidthSizing:   p.WidthSizing,
		HeightSizing:  p.HeightSizing,
		paddingTop:    padTop,
		paddingRight:  padRight,
		paddingBottom: padBottom,
		paddingLeft:   padLeft,
		gap:           p.Gap,
		alignment:     p.Align,
		justify:       p.Justify,
	}

	switch element.Type {
	case ElementText:
		w := 0
		for _, ch := range element.Text {
			w += RuneWidth(ch)
		}
		l.WidthSizing = Fixed(w)
		l.HeightSizing = Fixed(1)
	case ElementMultilineText:
		if element.Wrap {
			l.WidthSizing = Grow(1)
			l.HeightSizing = Fit()
			text := element.Text
			l.reflow = func(width int) int {
				if width <= 0 {
					return 1
				}
				return len(wrappedLines(text, width))
			}
		} else {
			lines := rawLines(element.Text)
			widest := 0
			for _, line := range lines {
				w := 0
				for _, ch := range line {
					w += RuneWidth(ch)
				}
				if w > widest {
					widest = w
				}
			}
			l.WidthSizing = Fixed(widest)
			l.HeightSizing = Fixed(len(lines))
		}
	}

	for _, child := range element.Children {
		l.Children = append(l.Children, buildLayoutTree(child))
	}
	return l
}

// paint walks the node tree in depth-first pre-order, matching the same
// traversal order that ComputeLayout uses to produce rects. parentStyle is
// the inherited style from ancestors; each element merges its own Style
// onto it and passes the result to children, so unspecified fields fall
// through the tree.
func paint(element Element, rects []Rect, idx int, screen *Screen, parentStyle Style) int {
	rect := rects[idx]
	idx++

	effective := mergeStyles(parentStyle, element.Style)

	switch element.Type {
	case ElementBox:
		for x := rect.X; x < rect.X+rect.Width; x++ {
			for y := rect.Y; y < rect.Y+rect.Height; y++ {
				screen.SetCell(x, y, ' ', effective)
			}
		}
		paintBorder(screen, rect, effective, element.Style.border)

	case ElementText:
		x := rect.X
		for _, ch := range element.Text {
			if x >= rect.X+rect.Width {
				break
			}
			screen.SetCell(x, rect.Y, ch, effective)
			x += RuneWidth(ch)
		}
	case ElementMultilineText:
		var lines []string
		if element.Wrap {
			lines = wrappedLines(element.Text, rect.Width)
		} else {
			lines = rawLines(element.Text)
		}
		for i, line := range lines {
			y := rect.Y + i
			if y >= rect.Y+rect.Height {
				break
			}
			x := rect.X
			for _, ch := range line {
				if x >= rect.X+rect.Width {
					break
				}
				screen.SetCell(x, y, ch, effective)
				x += RuneWidth(ch)
			}
		}
	}

	for _, child := range element.Children {
		idx = paint(child, rects, idx, screen, effective)
	}
	return idx
}

var Renderer = NewRenderer(StdOutScreen)

func paintBorder(screen *Screen, rect Rect, base Style, b Border) {
	if !b.Any() || rect.Width == 0 || rect.Height == 0 {
		return
	}

	bs := base
	if b.Color.Type != ColorNone {
		bs.foreground = b.Color
	}
	c := b.Chars

	x0, y0 := rect.X, rect.Y
	x1, y1 := rect.X+rect.Width-1, rect.Y+rect.Height-1

	if b.Top {
		for x := x0 + 1; x < x1; x++ {
			screen.SetCell(x, y0, c.Top, bs)
		}
	}
	if b.Bottom && y1 != y0 {
		for x := x0 + 1; x < x1; x++ {
			screen.SetCell(x, y1, c.Bottom, bs)
		}
	}
	if b.Left {
		for y := y0 + 1; y < y1; y++ {
			screen.SetCell(x0, y, c.Left, bs)
		}
	}
	if b.Right && x1 != x0 {
		for y := y0 + 1; y < y1; y++ {
			screen.SetCell(x1, y, c.Right, bs)
		}
	}

	if g := cornerGlyph(c.TopLeft, c.Top, c.Left, b.Top, b.Left); g != 0 {
		screen.SetCell(x0, y0, g, bs)
	}
	if g := cornerGlyph(c.TopRight, c.Top, c.Right, b.Top, b.Right); g != 0 {
		screen.SetCell(x1, y0, g, bs)
	}
	if g := cornerGlyph(c.BottomLeft, c.Bottom, c.Left, b.Bottom, b.Left); g != 0 {
		screen.SetCell(x0, y1, g, bs)
	}
	if g := cornerGlyph(c.BottomRight, c.Bottom, c.Right, b.Bottom, b.Right); g != 0 {
		screen.SetCell(x1, y1, g, bs)
	}
}

// cornerGlyph picks the rune to draw at a single corner of a box border.
// cornerChar is the full corner glyph (e.g. '┌' for top-left). hChar is the
// horizontal edge glyph (e.g. '─') and vChar is the vertical edge glyph
// (e.g. '│'). hasH and hasV report whether the horizontal/vertical edges
// meeting this corner are active. Return 0 to skip drawing the corner cell.
func cornerGlyph(cornerChar, hChar, vChar rune, hasH, hasV bool) rune {
	switch {
	case hasH && hasV:
		return cornerChar
	case hasH:
		return hChar
	case hasV:
		return vChar
	default:
		return 0
	}
}
