package tuix

const (
	ElementBox ElementType = iota
	ElementText
	ElementComponent
)

func Box(props Props, children ...Element) Element {
	return Element{
		Type: ElementBox,
		Layout: LayoutProps{
			Direction:     props.Direction,
			Gap:           props.Gap,
			PaddingTop:    props.Padding[0],
			PaddingRight:  props.Padding[1],
			PaddingBottom: props.Padding[2],
			PaddingLeft:   props.Padding[3],
		},
		Children: children,
	}
}

func Text(text string, style Style) Element {
	return Element{
		Type:  ElementText,
		Text:  text,
		Style: style,
	}
}

// Component wraps a render function and its props into an Element.
// The reconciler calls Render(props) to produce the concrete subtree.
func Component(fn func(Props) []Element, props Props) Element {
	return Element{
		Type: ElementComponent,
		Render: func(e Element) Element {
			children := fn(props)
			return Box(Props{}, children...)
		},
	}
}
