package tuix

import (
	"crypto/rand"
	"encoding/hex"
)

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
	Id       string
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

func createElementId() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CreateElement(children ...Element) Element {
	return Element{
		Id:       createElementId(),
		Children: children,
	}
}
