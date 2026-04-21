package tuix

func Reconcile(old *Node, next Element) (root *Node, deletions []*Node) {
	if old == nil {
		return buildTree(next, nil), nil
	}
	return diffNode(old, next, nil)
}

func buildTree(next Element, parent *Node) *Node {
	if next.Type == ElementComponent && next.Render != nil {
		next = next.Render(next)
	}
	node := createNode(next)
	node.Parent = parent
	if parent != nil {
		parent.Children = append(parent.Children, node)
	}
	return node
}

func diffChildren(oldChildren []*Node, next []Element, old *Node) ([]*Node, []*Node) {
	lut := make(map[string]*Node)
	var unkeyed []*Node

	var updatedNodes []*Node
	var deletionNodes []*Node

	for _, child := range oldChildren {
		if child.Element.Key != "" {
			lut[child.Element.Key] = child
			continue
		}
		unkeyed = append(unkeyed, child)
	}

	unkeyedIdx := 0

	for _, element := range next {

		if element.Key != "" {
			node := lut[element.Key]
			if node != nil {
				delete(lut, element.Key)
				updated, deletions := diffNode(node, element, node.Parent)
				updatedNodes = append(updatedNodes, updated)
				deletionNodes = append(deletionNodes, deletions...)
			} else {
				updated := buildTree(element, old)
				updatedNodes = append(updatedNodes, updated)
			}
		} else {
			if unkeyedIdx < len(unkeyed) {
				node := unkeyed[unkeyedIdx]
				unkeyedIdx++
				updated, deletions := diffNode(node, element, node.Parent)
				updatedNodes = append(updatedNodes, updated)
				deletionNodes = append(deletionNodes, deletions...)

			} else {
				updated := buildTree(element, old)
				updatedNodes = append(updatedNodes, updated)
			}
		}

	}

	for _, node := range unkeyed[unkeyedIdx:] {
		deletionNodes = append(deletionNodes, node)
	}
	for _, node := range lut {
		deletionNodes = append(deletionNodes, node)
	}

	return updatedNodes, deletionNodes
}

func diffNode(old *Node, next Element, parent *Node) (*Node, []*Node) {
	var deletions []*Node

	if old.Element.Type != next.Type || old.Element.Key != next.Key {
		// Type or key mismatch — replace entirely.
		deletions = append(deletions, old)
		return buildTree(next, parent), deletions
	}

	// Same type: update in place.
	if next.Type == ElementComponent && next.Render != nil {
		next = next.Render(next)
	}
	old.Element = next
	old.Parent = parent

	// Diff children.
	old.Children, deletions = diffChildren(old.Children, next.Children, old)
	return old, deletions
}
