package tuix

const (
	ElementBox ElementType = iota
	ElementText
	ElementMultilineText
	ElementComponent
)

func Box(props Props, style Style, children ...Element) Element {
	return Element{
		Type: ElementBox,
		Layout: LayoutProps{
			Direction:     props.Direction,
			WidthSizing:   Fit(),
			HeightSizing:  Fit(),
			Gap:           props.Gap,
			PaddingTop:    props.Padding[0],
			PaddingRight:  props.Padding[1],
			PaddingBottom: props.Padding[2],
			PaddingLeft:   props.Padding[3],
		},
		Style:    style,
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

// MultilineText renders text that may contain '\n' line breaks. Each '\n'
// starts a new row at the element's left edge; the intrinsic width is the
// longest line and the intrinsic height is the line count.
func MultilineText(text string, style Style) Element {
	return Element{
		Type:  ElementMultilineText,
		Text:  text,
		Style: style,
	}
}
