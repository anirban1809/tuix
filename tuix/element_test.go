package tuix

import "testing"

// ---- Box constructor -------------------------------------------------------

func TestBox_SetsHostElementType(t *testing.T) {
	el := Box(Props{})
	if el.Type != HostElement {
		t.Errorf("got Type=%d, want HostElement", el.Type)
	}
	if el.Tag != "box" {
		t.Errorf("got Tag=%q, want \"box\"", el.Tag)
	}
}

func TestBox_AttachesChildren(t *testing.T) {
	child1 := Text("a", Style{})
	child2 := Text("b", Style{})
	el := Box(Props{}, child1, child2)
	if len(el.Children) != 2 {
		t.Fatalf("got %d children, want 2", len(el.Children))
	}
	if el.Children[0].Props.Content != "a" {
		t.Errorf("child[0] content=%q, want \"a\"", el.Children[0].Props.Content)
	}
}

func TestBox_NoChildren_EmptySlice(t *testing.T) {
	el := Box(Props{})
	if len(el.Children) != 0 {
		t.Errorf("got %d children, want 0", len(el.Children))
	}
}

func TestBox_PropsPreserved(t *testing.T) {
	p := Props{
		Direction:    Column,
		WidthSizing:  Fixed(40),
		HeightSizing: Fixed(10),
		Gap:          2,
		Key:          "container",
	}
	el := Box(p)
	if el.Props.Direction != Column {
		t.Errorf("Direction not preserved")
	}
	if el.Props.Gap != 2 {
		t.Errorf("Gap not preserved")
	}
	if el.Props.Key != "container" {
		t.Errorf("Key not preserved")
	}
}

// ---- Text constructor -------------------------------------------------------

func TestText_SetsHostElementType(t *testing.T) {
	el := Text("hello", Style{})
	if el.Type != HostElement {
		t.Errorf("got Type=%d, want HostElement", el.Type)
	}
	if el.Tag != "text" {
		t.Errorf("got Tag=%q, want \"text\"", el.Tag)
	}
}

func TestText_ContentStoredInProps(t *testing.T) {
	el := Text("world", Style{})
	if el.Props.Content != "world" {
		t.Errorf("got Content=%q, want \"world\"", el.Props.Content)
	}
}

func TestText_StyleStoredInProps(t *testing.T) {
	s := Style{}.Bold(true)
	el := Text("hi", s)
	if !el.Props.Style.IsBold() {
		t.Errorf("style not stored in props")
	}
}

func TestText_DefaultSizing(t *testing.T) {
	el := Text("hi", Style{})
	if el.Props.WidthSizing.Mode != SizingFit {
		t.Errorf("WidthSizing: got %d, want SizingFit", el.Props.WidthSizing.Mode)
	}
	if el.Props.HeightSizing.Mode != SizingFixed || el.Props.HeightSizing.Value != 1 {
		t.Errorf("HeightSizing: got {%d %d}, want Fixed(1)", el.Props.HeightSizing.Mode, el.Props.HeightSizing.Value)
	}
}

func TestText_HasNoChildren(t *testing.T) {
	el := Text("hi", Style{})
	if len(el.Children) != 0 {
		t.Errorf("text element should have no children, got %d", len(el.Children))
	}
}

// ---- Component constructor --------------------------------------------------

func TestComponent_SetsFuncElementType(t *testing.T) {
	fn := func(p Props) []Element { return nil }
	el := Component(fn, Props{})
	if el.Type != FuncElement {
		t.Errorf("got Type=%d, want FuncElement", el.Type)
	}
}

func TestComponent_RenderFnSet(t *testing.T) {
	fn := func(p Props) []Element { return nil }
	el := Component(fn, Props{})
	if el.RenderFn == nil {
		t.Errorf("RenderFn should not be nil")
	}
}

func TestComponent_RenderFnExecutable(t *testing.T) {
	called := false
	fn := func(p Props) []Element {
		called = true
		return []Element{Text("x", Style{})}
	}
	el := Component(fn, Props{})
	result := el.RenderFn(el.Props)
	if !called {
		t.Errorf("RenderFn was not called")
	}
	if len(result) != 1 {
		t.Errorf("got %d children, want 1", len(result))
	}
}

func TestComponent_PropsForwardedToRenderFn(t *testing.T) {
	var receivedKey string
	fn := func(p Props) []Element {
		receivedKey = p.Key
		return nil
	}
	el := Component(fn, Props{Key: "myKey"})
	el.RenderFn(el.Props)
	if receivedKey != "myKey" {
		t.Errorf("props.Key not forwarded: got %q, want \"myKey\"", receivedKey)
	}
}

// ---- Key field ---------------------------------------------------------------

func TestElement_KeyIsPreservedOnBox(t *testing.T) {
	el := Box(Props{Key: "k1"})
	if el.Props.Key != "k1" {
		t.Errorf("Key=%q, want \"k1\"", el.Props.Key)
	}
}
