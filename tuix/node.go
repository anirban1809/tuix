package tuix

type ElementType int

type Props struct {
	Direction Direction
	Gap       int
	Padding   [4]int
	Values    map[string]any
}

func (p Props) Get(key string) any {
	return p.Values[key]
}

type Element struct {
	Type     ElementType
	Key      string
	Text     string
	Style    Style
	Layout   LayoutProps
	Children []Element
	Render   func(props Element) Element
	Props    Props
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

type Node struct {
	Element   Element
	Rect      Rect
	Children  []*Node
	Parent    *Node
	HookState HookState
}

func createNode(element Element) *Node {
	var children []*Node
	node := &Node{
		HookState: HookState{},
	}

	for _, child := range element.Children {
		childNode := createNode(child)
		childNode.Parent = node
		children = append(children, childNode)
	}

	node.Element = element
	node.Children = children

	return node
}
