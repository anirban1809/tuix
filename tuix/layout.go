package tuix

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Direction int

const (
	Row Direction = iota
	Column
)

type SizingMode int

const (
	SizingFixed SizingMode = iota
	SizingGrow
	SizingFit
)

type Sizing struct {
	Mode  SizingMode
	Value int
}

type LayoutNode struct {
	Direction       Direction
	WidthSizing     Sizing
	HeightSizing    Sizing
	Children        []*LayoutNode
	intrinsicHeight int
	intrinsicWidth  int
	paddingTop      int
	paddingBottom   int
	paddingLeft     int
	paddingRight    int
	gap             int
	alignment       Alignment
	justify         Justify
	// reflow lets a node recompute its main-axis size once its cross-axis
	// size is known during the layout pass. Used by content whose height
	// depends on width (e.g. word-wrapped text). Currently invoked only when
	// the parent's Direction is Column.
	reflow func(crossSize int) int
}

func NewLayout() *LayoutNode {
	return &LayoutNode{}
}

func (l LayoutNode) WithDirection(dir Direction) *LayoutNode {
	l.Direction = dir
	return &l
}

func Fixed(n int) Sizing { return Sizing{Mode: SizingFixed, Value: n} }
func Grow(n int) Sizing  { return Sizing{Mode: SizingGrow, Value: n} }
func Fit() Sizing        { return Sizing{Mode: SizingFit} }

func (l *LayoutNode) WithSize(w, h Sizing) *LayoutNode {
	l.WidthSizing = w
	l.HeightSizing = h
	return l
}

func (l *LayoutNode) WithChildren(children ...*LayoutNode) *LayoutNode {
	l.Children = children
	return l
}

func (l *LayoutNode) WithPadding(top, right, bottom, left int) *LayoutNode {
	l.paddingTop, l.paddingRight, l.paddingBottom, l.paddingLeft = top, right, bottom, left
	return l
}

func (l *LayoutNode) WithGap(value int) *LayoutNode {
	l.gap = value
	return l
}

type Alignment int

const (
	AlignStretch Alignment = iota
	AlignStart
	AlignCenter
	AlignEnd
)

func (l *LayoutNode) WithAlign(alignment Alignment) *LayoutNode {
	l.alignment = alignment
	return l
}

type Justify int

const (
	JustifyStart Justify = iota
	JustifyEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
)

func (l *LayoutNode) WithJustify(value Justify) *LayoutNode {
	l.justify = value
	return l
}
