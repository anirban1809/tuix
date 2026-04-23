package tuix

type ComponentRenderer struct {
	screen *Screen
	dirty  chan struct{}
}

func NewRenderer(screen *Screen) *ComponentRenderer {
	return &ComponentRenderer{screen: screen, dirty: make(chan struct{})}
}

func (r *ComponentRenderer) Render(next Element) {
	available := Rect{X: 0, Y: 0, Width: r.screen.Width(), Height: r.screen.Height()}
	layoutRoot := buildLayoutTree(next)
	rects := ComputeLayout(layoutRoot, available)

	r.screen.Clear()
	paint(next, rects, 0, r.screen)
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

	if element.Text != "" {
		w := 0
		for _, ch := range element.Text {
			w += RuneWidth(ch)
		}
		l.WidthSizing = Fixed(w)
		l.HeightSizing = Fixed(1)
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

	if element.Text != "" {
		x := rect.X
		for _, ch := range element.Text {
			if x >= rect.X+rect.Width {
				break
			}
			screen.SetCell(x, rect.Y, ch, element.Style)
			x += RuneWidth(ch)
		}
	}

	for _, child := range element.Children {
		idx = paint(child, rects, idx, screen)
	}
	return idx
}

var Renderer = NewRenderer(StdOutScreen)
