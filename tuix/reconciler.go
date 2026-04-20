package tuix

// Reconciler diffs a new Element tree against the existing committed Fiber tree
// and produces an updated Fiber tree ready for the commit phase.
//
// Render cycle:
//  1. ScheduleUpdate builds a work-in-progress root fiber.
//  2. performUnitOfWork processes one fiber at a time (depth-first).
//  3. reconcileChildren diffs new elements against existing fiber children.
//  4. commitRoot walks the completed work-in-progress tree and applies effects.
type Reconciler struct {
	// current is the root of the last committed fiber tree (what is on screen).
	current *Fiber

	// workInProgress is the root of the fiber tree being built this cycle.
	workInProgress *Fiber

	// deletions accumulates fibers that must be removed during commit.
	deletions []*Fiber
}

// NewReconciler creates a reconciler with an empty committed tree.
func NewReconciler() *Reconciler {
	return &Reconciler{}
}

// Reconcile runs a full render cycle for the given root element and returns
// the committed fiber tree root. The caller (Renderer) uses the fiber tree
// to drive layout and screen painting.
func (r *Reconciler) Reconcile(root Element) *Fiber {
	r.deletions = nil

	r.workInProgress = &Fiber{
		Tag:       FiberRoot,
		Alternate: r.current,
	}

	// Treat the root element as the single child of the wip root.
	r.reconcileChildren(r.workInProgress, []Element{root})

	// Walk the work-in-progress tree depth-first, completing each fiber.
	cursor := r.workInProgress.Child
	for cursor != nil {
		cursor = r.performUnitOfWork(cursor)
	}

	r.commitRoot()
	return r.current
}

// performUnitOfWork completes work for one fiber (calls render for FuncElement,
// reconciles its children) and returns the next fiber to process.
//
// Traversal order: go deep (child) first, then across (sibling), then up and
// across (return.sibling). This is the standard fiber work-loop pattern.
func (r *Reconciler) performUnitOfWork(f *Fiber) *Fiber {
	// Expand FuncElement fibers by calling their render function.
	if f.Tag == FiberFunc && f.RenderFn != nil {
		children := f.RenderFn(f.Props)
		r.reconcileChildren(f, children)
	} else if f.Tag == FiberHost {
		// HostElement children come from the element tree, already set during
		// reconcileChildren of the parent — re-reconcile to build fiber children.
		// (No-op if reconcileChildren was already called from the parent pass.)
	}

	// Depth-first: go to child first.
	if f.Child != nil {
		return f.Child
	}

	// No child: move to sibling or climb until we find an uncle.
	cur := f
	for cur != nil {
		if cur.Sibling != nil {
			return cur.Sibling
		}
		cur = cur.Return
	}
	return nil
}

// reconcileChildren diffs elements against the existing children of parent
// (found via parent.Alternate.Child) and attaches new/updated/deleted fiber
// nodes onto parent.Child and the sibling chain.
func (r *Reconciler) reconcileChildren(parent *Fiber, elements []Element) {
	// oldFiber is the head of the previous render's child list for this parent.
	var oldFiber *Fiber
	if parent.Alternate != nil {
		oldFiber = parent.Alternate.Child
	}

	// Build a key → oldFiber map for O(1) key-based matching.
	keyMap := make(map[string]*Fiber)
	for f := oldFiber; f != nil; f = f.Sibling {
		if f.ElementKey != "" {
			keyMap[f.ElementKey] = f
		}
	}

	// positionCursor walks the old fiber list for position-based (keyless) matching.
	positionCursor := oldFiber

	var prevSibling *Fiber

	for i, el := range elements {
		var matchedOld *Fiber

		key := el.Props.Key

		if key != "" {
			// Key-based match: find the old fiber with the same key.
			if candidate, ok := keyMap[key]; ok {
				matchedOld = candidate
				delete(keyMap, key) // consumed
			}
		} else {
			// Position-based match: reuse the old fiber at the same index if it
			// has the same tag (host "box"/"text" vs func).
			if positionCursor != nil && positionCursor.ElementKey == "" &&
				sameType(positionCursor, el) {
				matchedOld = positionCursor
				positionCursor = positionCursor.Sibling
			} else if positionCursor != nil {
				positionCursor = positionCursor.Sibling
			}
		}

		newFiber := r.createOrUpdateFiber(matchedOld, el, parent)

		// Wire siblings and child pointer.
		if i == 0 {
			parent.Child = newFiber
		} else if prevSibling != nil {
			prevSibling.Sibling = newFiber
		}
		prevSibling = newFiber
	}

	// Any old fibers that were not matched (by key or position) are deletions.
	// Key-map survivors are unmatched key fibers.
	for _, old := range keyMap {
		r.scheduleDeletion(old)
	}
	// Remaining position-cursor chain are unmatched position fibers.
	for f := positionCursor; f != nil; f = f.Sibling {
		if f.ElementKey == "" { // skip keyed ones already handled above
			r.scheduleDeletion(f)
		}
	}
}

// createOrUpdateFiber either reuses an existing fiber (EffectUpdate) or creates
// a new one (EffectPlace), wiring Alternate and Return appropriately.
func (r *Reconciler) createOrUpdateFiber(old *Fiber, el Element, parent *Fiber) *Fiber {
	if old != nil {
		// Reuse: update props, point Alternate to old fiber.
		old.Alternate.Props = old.Props // preserve old props in alternate
		updated := &Fiber{
			Tag:         fiberTagFor(el),
			ElementType: el.Type,
			ElementTag:  el.Tag,
			ElementKey:  el.Props.Key,
			Props:       el.Props,
			RenderFn:    el.RenderFn,
			Return:      parent,
			Alternate:   old,
			Effect:      EffectUpdate,
			LayoutNode:  old.LayoutNode, // carry forward; renderer will patch
		}
		// Eagerly reconcile host children from the element tree.
		if el.Type == HostElement && len(el.Children) > 0 {
			r.reconcileChildren(updated, el.Children)
		}
		return updated
	}

	// New fiber.
	placed := &Fiber{
		Tag:         fiberTagFor(el),
		ElementType: el.Type,
		ElementTag:  el.Tag,
		ElementKey:  el.Props.Key,
		Props:       el.Props,
		RenderFn:    el.RenderFn,
		Return:      parent,
		Effect:      EffectPlace,
	}
	if el.Type == HostElement && len(el.Children) > 0 {
		r.reconcileChildren(placed, el.Children)
	}
	return placed
}

// scheduleDeletion adds a fiber to the deletion list with EffectDelete.
func (r *Reconciler) scheduleDeletion(f *Fiber) {
	f.Effect = EffectDelete
	f.deletionNext = nil
	if len(r.deletions) > 0 {
		r.deletions[len(r.deletions)-1].deletionNext = f
	}
	r.deletions = append(r.deletions, f)
}

// commitRoot swaps work-in-progress → current and processes all effects.
func (r *Reconciler) commitRoot() {
	// Process deletions first.
	for _, f := range r.deletions {
		r.commitWork(f)
	}

	// Walk the new tree and commit Place/Update effects.
	walk(r.workInProgress, func(f *Fiber) {
		r.commitWork(f)
	})

	// Swap: the work-in-progress tree becomes the committed current tree.
	r.current = r.workInProgress
	r.workInProgress = nil
}

// commitWork applies the effect of a single fiber.
// In Phase 3 this is intentionally minimal: the renderer reads the fiber tree
// directly, so commit just ensures LayoutNodes exist for host fibers.
func (r *Reconciler) commitWork(f *Fiber) {
	switch f.Effect {
	case EffectPlace:
		if f.Tag == FiberHost {
			f.LayoutNode = fiberToLayoutNode(f)
		}
	case EffectUpdate:
		if f.Tag == FiberHost && f.LayoutNode != nil {
			applyPropsToLayoutNode(f.LayoutNode, f.Props)
		} else if f.Tag == FiberHost {
			f.LayoutNode = fiberToLayoutNode(f)
		}
	case EffectDelete:
		f.LayoutNode = nil
	}
}

// fiberToLayoutNode creates a LayoutNode from a FiberHost's props.
func fiberToLayoutNode(f *Fiber) *LayoutNode {
	n := NewLayout()
	applyPropsToLayoutNode(n, f.Props)
	return n
}

// applyPropsToLayoutNode syncs layout-relevant props onto an existing node.
func applyPropsToLayoutNode(n *LayoutNode, p Props) {
	n.Direction = p.Direction
	n.WidthSizing = p.WidthSizing
	n.HeightSizing = p.HeightSizing
	n.WithPadding(p.Padding[0], p.Padding[1], p.Padding[2], p.Padding[3])
	n.WithGap(p.Gap)
	n.WithAlign(p.Align)
	n.WithJustify(p.Justify)
}

// fiberTagFor maps an Element to the correct FiberTag.
func fiberTagFor(el Element) FiberTag {
	if el.Type == FuncElement {
		return FiberFunc
	}
	return FiberHost
}

// sameType returns true when an existing fiber and a new element are
// compatible enough for position-based reuse (same host tag or both funcs).
func sameType(f *Fiber, el Element) bool {
	if f.ElementType != el.Type {
		return false
	}
	if el.Type == HostElement {
		return f.ElementTag == el.Tag
	}
	// FuncElement: reuse if same function pointer (comparable in Go)
	return f.RenderFn != nil && el.RenderFn != nil
}
