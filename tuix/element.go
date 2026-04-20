package tuix

// ElementType identifies whether a node is a host primitive or a user-defined
// component function.
type ElementType int

const (
	HostElement ElementType = iota // maps directly to a LayoutNode + rendered content
	FuncElement                    // a user-defined component function
)

// Props carries all configuration for an element: layout, visual, and identity.
type Props struct {
	// Layout
	Direction    Direction
	WidthSizing  Sizing
	HeightSizing Sizing
	Padding      [4]int // top, right, bottom, left
	Gap          int
	Align        Alignment
	Justify      Justify

	// Visual
	Content string
	Style   Style

	// Identity — used by the reconciler for stable matching across re-renders
	Key string
}

// Element is a lightweight, stateless description of a node in the UI tree.
// Elements are recreated on every render; Fibers are the persistent counterpart.
type Element struct {
	Type     ElementType
	Tag      string    // "box" or "text" for HostElement; "" for FuncElement
	Props    Props
	Children []Element

	// RenderFn is set for FuncElement; nil for HostElement
	RenderFn ComponentFunc
}

// ComponentFunc is the signature all user-defined components must implement.
type ComponentFunc func(props Props) []Element

// Box creates a host element that renders as a layout container.
func Box(props Props, children ...Element) Element {
	return Element{
		Type:     HostElement,
		Tag:      "box",
		Props:    props,
		Children: children,
	}
}

// Text creates a host element that renders a string of content.
func Text(content string, style Style) Element {
	p := Props{
		Content:      content,
		Style:        style,
		WidthSizing:  Fit(),
		HeightSizing: Fixed(1),
	}
	return Element{
		Type:  HostElement,
		Tag:   "text",
		Props: p,
	}
}

// Component wraps a ComponentFunc as an Element so it can appear in the tree.
func Component(fn ComponentFunc, props Props, children ...Element) Element {
	return Element{
		Type:     FuncElement,
		Props:    props,
		Children: children,
		RenderFn: fn,
	}
}
