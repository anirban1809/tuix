package tuix

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
	paint(next, rects, 0, r.screen)

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
		widestLineLength := 0
		lineLength := 0
		height := 0

		for _, ch := range element.Text {
			if ch == '\n' {
				if lineLength > widestLineLength {
					widestLineLength = lineLength
				}
				lineLength = 0
				height++
			} else {
				lineLength += RuneWidth(ch)
			}
		}

		l.WidthSizing = Fixed(widestLineLength)
		l.HeightSizing = Fixed(height)
	}

	for _, child := range element.Children {
		l.Children = append(l.Children, buildLayoutTree(child))
	}
	return l
}

// paint walks the node tree in depth-first pre-order, matching the same
// traversal order that ComputeLayout uses to produce rects.
func paint(element Element, rects []Rect, idx int, screen *Screen) int {
	rect := rects[idx]
	idx++

	switch element.Type {
	case ElementText:
		x := rect.X
		for _, ch := range element.Text {
			if x >= rect.X+rect.Width {
				break
			}
			screen.SetCell(x, rect.Y, ch, element.Style)
			x += RuneWidth(ch)
		}
	case ElementMultilineText:
		x := rect.X
		y := rect.Y
		for _, ch := range element.Text {
			if ch == '\n' {
				y++
				x = rect.X
				continue
			}
			if y >= rect.Y+rect.Height {
				break
			}
			if x < rect.X+rect.Width {
				screen.SetCell(x, y, ch, element.Style)
			}
			x += RuneWidth(ch)
		}
	}

	for _, child := range element.Children {
		idx = paint(child, rects, idx, screen)
	}
	return idx
}

var Renderer = NewRenderer(StdOutScreen)
