package suffixtree

type edge struct {
	label []rune
	*node
}

func newEdge(label []rune, node *node) *edge {
	return &edge{label: label, node: node}
}
