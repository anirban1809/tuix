package tuix

import "strings"

// Renderer ties the reconciler, layout engine, and screen together.
//
// Each call to Render runs the full pipeline:
//  1. Reconciler diffs elements → committed fiber tree
//  2. buildLayoutTree converts host fibers → LayoutNode tree
//  3. ComputeLayout computes []Rect (depth-first pre-order)
//  4. paintFibers walks fibers + rects → Screen.SetCell
//  5. Screen.Flush emits only changed ANSI sequences
type Renderer struct {
	screen     *Screen
	reconciler *Reconciler
	width      int
	height     int
}

// NewRenderer creates a renderer that paints into the given screen.
func NewRenderer(screen *Screen) *Renderer {
	return &Renderer{
		screen:     screen,
		reconciler: NewReconciler(),
		width:      screen.Width(),
		height:     screen.Height(),
	}
}

// Render runs a full render cycle for the given root element.
func (r *Renderer) Render(root Element) {
	committedRoot := r.reconciler.Reconcile(root)
	if committedRoot == nil {
		return
	}

	// Build the LayoutNode tree from the committed fiber tree.
	available := Rect{X: 0, Y: 0, Width: r.width, Height: r.height}
	layoutRoot := r.buildLayoutTree(committedRoot.Child)
	if layoutRoot == nil {
		return
	}

	rects := ComputeLayout(layoutRoot, available)

	// Walk host fibers in the same depth-first pre-order as ComputeLayout and
	// paint each one using its corresponding rect.
	idx := 0
	r.paintFibers(committedRoot.Child, rects, &idx)

	r.screen.Flush()
}

// buildLayoutTree converts the host-fiber subtree rooted at f into a LayoutNode
// tree. FuncElement fibers are transparent — we descend into their children.
// Returns nil if there are no host fibers.
func (r *Renderer) buildLayoutTree(f *Fiber) *LayoutNode {
	if f == nil {
		return nil
	}

	// Skip func fibers transparently.
	if f.Tag == FiberFunc {
		return r.buildLayoutTree(f.Child)
	}

	node := f.LayoutNode
	if node == nil {
		return nil
	}

	// Recursively attach child layout nodes.
	node.Children = nil
	child := f.Child
	for child != nil {
		childNode := r.buildLayoutTree(child)
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
		child = child.Sibling
	}

	return node
}

// paintFibers walks host fibers in depth-first pre-order (matching ComputeLayout
// output order) and writes cell content into the screen buffer.
func (r *Renderer) paintFibers(f *Fiber, rects []Rect, idx *int) {
	if f == nil {
		return
	}

	if f.Tag == FiberHost {
		if *idx < len(rects) {
			rect := rects[*idx]
			*idx++

			if f.ElementTag == "text" {
				r.paintText(f.Props.Content, f.Props.Style, rect)
			}
			// "box" host elements have no direct content; children are painted below.
		}
	}

	// Recurse into children, then siblings.
	r.paintFibers(f.Child, rects, idx)
	r.paintFibers(f.Sibling, rects, idx)
}

// paintText writes the content string into the screen starting at rect's origin,
// clipping to rect's width.
func (r *Renderer) paintText(content string, style Style, rect Rect) {
	runes := []rune(strings.TrimRight(content, "\n"))
	x := rect.X
	for _, ch := range runes {
		if x >= rect.X+rect.Width {
			break
		}
		r.screen.SetCell(x, rect.Y, ch, style)
		x++
	}
}
