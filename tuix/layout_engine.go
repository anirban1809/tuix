package tuix

// IntrinsicSize returns the natural (unconstrained) dimensions of a layout tree.
func IntrinsicSize(root *LayoutNode) (width, height int) {
	return measure(root)
}

// ComputeLayout runs the two-pass layout algorithm on the given tree and
// returns a flat slice of Rects in depth-first, pre-order (parent before
// children). The root is given the provided available Rect as its bounds.
func ComputeLayout(root *LayoutNode, available Rect) []Rect {
	// Pass 1 — measure intrinsic sizes bottom-up.
	root.intrinsicWidth, root.intrinsicHeight = measure(root)

	// Pass 2 — assign concrete rects top-down.
	var out []Rect
	layout(root, available, &out)
	return out
}

// measure fills in intrinsicWidth and intrinsicHeight for every node in the
// subtree rooted at n (bottom-up, children first).
func measure(n *LayoutNode) (int, int) {
	for _, child := range n.Children {
		child.intrinsicWidth, child.intrinsicHeight = measure(child)
	}

	width, height := 0, 0

	gaps := 0
	if len(n.Children) > 1 {
		gaps = n.gap * (len(n.Children) - 1)
	}

	switch n.WidthSizing.Mode {
	case SizingFixed:
		width = n.WidthSizing.Value

	case SizingFit:
		if n.Direction == Row {
			for _, child := range n.Children {
				width += child.intrinsicWidth
			}
			width += gaps
		} else {
			for _, child := range n.Children {
				if child.intrinsicWidth > width {
					width = child.intrinsicWidth
				}
			}
		}

	case SizingGrow:
		width = 0
	}

	switch n.HeightSizing.Mode {
	case SizingFixed:
		height = n.HeightSizing.Value

	case SizingFit:
		if n.Direction != Row {
			for _, child := range n.Children {
				height += child.intrinsicHeight
			}
			height += gaps
		} else {
			for _, child := range n.Children {
				if child.intrinsicHeight > height {
					height = child.intrinsicHeight
				}
			}
		}

	case SizingGrow:
		height = 0

	}

	if n.WidthSizing.Mode == SizingFit {
		width += n.paddingLeft + n.paddingRight
	}
	if n.HeightSizing.Mode == SizingFit {
		height += n.paddingTop + n.paddingBottom
	}

	return width, height
}

func mainSize(r Rect, dir Direction) int {
	if dir == Row {
		return r.Width
	}
	return r.Height
}

func setMainSize(r *Rect, dir Direction, value int) {
	if dir == Row {
		r.Width = value
	} else {
		r.Height = value
	}
}

func setCrossSize(r *Rect, dir Direction, value int) {
	if dir == Row {
		r.Height = value
	} else {
		r.Width = value
	}
}

func mainStart(r Rect, dir Direction) int {
	if dir == Row {
		return r.X
	}
	return r.Y
}

func crossAvailableSize(r Rect, dir Direction) int {
	if dir == Row {
		return r.Height
	}
	return r.Width
}

func resolveCrossSize(parent *LayoutNode, child *LayoutNode, into Rect) int {
	if parent.alignment == AlignStretch {
		return crossAvailableSize(into, parent.Direction)
	}

	if parent.Direction == Row {
		switch child.HeightSizing.Mode {
		case SizingFixed, SizingFit:
			return child.intrinsicHeight
		case SizingGrow:
			return into.Height
		}
	}

	switch child.WidthSizing.Mode {
	case SizingFixed, SizingFit:
		return child.intrinsicWidth
	case SizingGrow:
		return into.Width
	}

	return 0
}

func applyCrossAlignment(parent *LayoutNode, childRect *Rect, into Rect) {
	if parent.Direction == Row {
		switch parent.alignment {
		case AlignCenter:
			childRect.Y = into.Y + (into.Height-childRect.Height)/2
		case AlignEnd:
			childRect.Y = into.Y + (into.Height - childRect.Height)
		default:
			childRect.Y = into.Y
		}
		return
	}

	switch parent.alignment {
	case AlignCenter:
		childRect.X = into.X + (into.Width-childRect.Width)/2
	case AlignEnd:
		childRect.X = into.X + (into.Width - childRect.Width)
	default:
		childRect.X = into.X
	}
}

func resolveJustify(justify Justify, start, innerMainSize, usedMain, baseGap, childCount int) (int, int) {
	if childCount == 0 {
		return start, 0
	}

	minGapTotal := 0
	if childCount > 1 {
		minGapTotal = baseGap * (childCount - 1)
	}

	extraSpace := innerMainSize - usedMain - minGapTotal
	if extraSpace < 0 {
		extraSpace = 0
	}

	cursor := start
	gap := baseGap

	switch justify {
	case JustifyEnd:
		cursor += extraSpace
	case JustifyCenter:
		cursor += extraSpace / 2
	case JustifySpaceBetween:
		if childCount > 1 {
			gap += extraSpace / (childCount - 1)
		}
	case JustifySpaceAround:
		segment := extraSpace / (childCount * 2)
		cursor += segment
		gap += segment * 2
	}

	return cursor, gap
}

// layout assigns a concrete Rect to n given the space offered by its parent,
// then recurses into children (top-down).
func layout(n *LayoutNode, into Rect, out *[]Rect) {
	*out = append(*out, into)

	into.X += n.paddingLeft
	into.Y += n.paddingTop
	into.Width -= n.paddingLeft + n.paddingRight
	into.Height -= n.paddingTop + n.paddingBottom

	if into.Width < 0 {
		into.Width = 0
	}
	if into.Height < 0 {
		into.Height = 0
	}

	if len(n.Children) == 0 {
		return
	}

	innerMainSize := mainSize(into, n.Direction)
	totalGap := n.gap * (len(n.Children) - 1)
	usedByFixedAndFit := 0
	totalGrowWeight := 0
	growIndices := make([]int, 0, len(n.Children))
	childRects := make([]Rect, 0, len(n.Children))

	for _, child := range n.Children {
		var childRect Rect

		setCrossSize(&childRect, n.Direction, resolveCrossSize(n, child, into))

		if n.Direction == Row {
			switch child.WidthSizing.Mode {
			case SizingFixed, SizingFit:
				childRect.Width = child.intrinsicWidth
				usedByFixedAndFit += childRect.Width
			case SizingGrow:
				totalGrowWeight += child.WidthSizing.Value
				growIndices = append(growIndices, len(childRects))
			}
		} else {
			switch child.HeightSizing.Mode {
			case SizingFixed, SizingFit:
				childRect.Height = child.intrinsicHeight
				usedByFixedAndFit += childRect.Height
			case SizingGrow:
				totalGrowWeight += child.HeightSizing.Value
				growIndices = append(growIndices, len(childRects))
			}
		}

		childRects = append(childRects, childRect)
	}

	remaining := max(innerMainSize-totalGap-usedByFixedAndFit, 0)

	if totalGrowWeight > 0 {
		remainingWeight := totalGrowWeight
		remainingSpace := remaining

		for _, idx := range growIndices {
			child := n.Children[idx]
			weight := child.WidthSizing.Value
			if n.Direction == Column {
				weight = child.HeightSizing.Value
			}

			size := remainingSpace
			if remainingWeight > weight {
				size = remainingSpace * weight / remainingWeight
			}

			setMainSize(&childRects[idx], n.Direction, size)
			remainingSpace -= size
			remainingWeight -= weight
		}
	}

	usedMain := 0
	for _, childRect := range childRects {
		usedMain += mainSize(childRect, n.Direction)
	}

	cursor, gap := resolveJustify(
		n.justify,
		mainStart(into, n.Direction),
		innerMainSize,
		usedMain,
		n.gap,
		len(n.Children),
	)

	for i, child := range n.Children {
		childRect := childRects[i]

		if n.Direction == Row {
			childRect.X = cursor
		} else {
			childRect.Y = cursor
		}

		applyCrossAlignment(n, &childRect, into)
		layout(child, childRect, out)

		cursor += mainSize(childRect, n.Direction)
		if i < len(n.Children)-1 {
			cursor += gap
		}
	}
}
