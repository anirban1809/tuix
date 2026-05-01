package tuix

import (
	"strings"
)

// multilineLines returns the rendered lines for an ElementMultilineText. It
// always splits on '\n' first; if WrapWidth > 0, each segment is further
// broken at word boundaries by wrapText so it fits within WrapWidth columns.
func multilineLines(element Element) []string {
	segments := strings.Split(element.Text, "\n")
	if element.WrapWidth <= 0 {
		return segments
	}
	var out []string
	for _, seg := range segments {
		out = append(out, wrapText(seg, element.WrapWidth)...)
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
		if i%maxWidth == 0 {
			lines = append(lines, line.String())
			line.Reset()
			continue
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
	l := &LayoutNode{
		Direction:     p.Direction,
		WidthSizing:   p.WidthSizing,
		HeightSizing:  p.HeightSizing,
		paddingTop:    p.PaddingTop,
		paddingRight:  p.PaddingRight,
		paddingBottom: p.PaddingBottom,
		paddingLeft:   p.PaddingLeft,
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
		lines := multilineLines(element)
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
		lines := multilineLines(element)
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
