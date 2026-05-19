package tuix

// ElementMarkdown is a marker type used to distinguish markdown content
// from regular text elements. It carries pre-parsed markdown cells that
// include per-character styling (colors, bold, italic, etc.).
type MarkdownContent struct {
    Lines []markdownLine
}

type ElementType int

type Props struct {
	Direction Direction
	Gap       int
	Padding   [4]int
	Align     Alignment
	Justify   Justify
	// Width/Height sizing. The zero value (Sizing{} == Fixed(0)) is treated
	// as unset and defaults to Fit() inside Box. Use Grow(1) to fill the
	// parent's cross axis (or the terminal, when applied to the root).
	Width  Sizing
	Height Sizing
	Values map[string]any
}

func (p Props) Get(key string) any {
	return p.Values[key]
}

type Element struct {
	Id       string
	Type     ElementType
	Key      string
	Text     string
	Wrap     bool
	Style    Style
	Layout   LayoutProps
	Children []Element
	Render   func(props Element) Element
	Props    Props
	Markdown     MarkdownContent
	MarkdownText string // raw markdown string for re-parsing during layout
}

type LayoutProps struct {
	Direction     Direction
	WidthSizing   Sizing
	HeightSizing  Sizing
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int
	PaddingLeft   int
	Gap           int
	Align         Alignment
	Justify       Justify
}
