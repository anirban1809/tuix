package tuix

// FiberTag classifies the kind of work a fiber represents.
type FiberTag int

const (
	FiberRoot FiberTag = iota // invisible root, anchors the tree
	FiberHost                 // HostElement — produces a LayoutNode
	FiberFunc                 // FuncElement — calls a ComponentFunc
)

// EffectTag describes what commit-phase action is needed for a fiber.
type EffectTag int

const (
	EffectNone   EffectTag = iota
	EffectPlace            // new fiber, needs insertion into the output tree
	EffectUpdate           // props changed, needs patch
	EffectDelete           // removed from the new tree, needs cleanup
)

// Fiber is a long-lived work unit that persists across render cycles.
// The tree uses three pointers (Child, Sibling, Return) so the entire tree
// can be traversed iteratively with a single cursor — no recursion needed.
type Fiber struct {
	Tag FiberTag

	// Element data this fiber was built from
	ElementType ElementType
	ElementTag  string // "box", "text", or "" for FuncElement
	ElementKey  string // from Props.Key; drives key-based reconciliation

	// Tree structure
	Child   *Fiber // first child
	Sibling *Fiber // next sibling (links all children as a linked list)
	Return  *Fiber // parent fiber

	// Alternate points to the fiber from the previous committed render that
	// this fiber corresponds to. current.Alternate = workInProgress, and vice
	// versa after commit — this is the double-buffer pattern for the fiber tree.
	Alternate *Fiber

	// Props from the most recent render of this fiber
	Props Props

	// RenderFn is set for FiberFunc fibers; nil for FiberHost
	RenderFn ComponentFunc

	// LayoutNode is the concrete layout node produced by FiberHost fibers.
	// FiberFunc fibers leave this nil; the bridge walks down to host fibers.
	LayoutNode *LayoutNode

	// Effect describes what the commit phase should do with this fiber.
	Effect EffectTag

	// deletionNext links fibers in the reconciler's pending-deletion list.
	deletionNext *Fiber
}

// firstChild is a nil-safe helper that returns the first child of the fiber.
func (f *Fiber) firstChild() *Fiber {
	if f == nil {
		return nil
	}
	return f.Child
}

// nextSibling is a nil-safe helper that returns the next sibling.
func (f *Fiber) nextSibling() *Fiber {
	if f == nil {
		return nil
	}
	return f.Sibling
}

// walk visits every fiber in the subtree rooted at f in depth-first pre-order,
// calling fn for each. Used by tests and the renderer to traverse the tree
// without managing the three-pointer walk manually.
func walk(f *Fiber, fn func(*Fiber)) {
	if f == nil {
		return
	}
	fn(f)
	child := f.Child
	for child != nil {
		walk(child, fn)
		child = child.Sibling
	}
}
