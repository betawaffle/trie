package trie

type copyTxn struct {
	nodes []Node
	edges []*Node
}

func (root *Node) DenseCopy() *Node {
	var n int
	root.Walk(func(node *Node) bool {
		n += 1
		return true
	})
	t := &copyTxn{
		nodes: make([]Node, n),
		edges: make([]*Node, n-1),
	}
	return t.copyNode(root)
}

func (t *copyTxn) copyNode(n *Node) (cp *Node) {
	if i := len(t.nodes) - 1; i >= 0 {
		t.nodes, cp = t.nodes[:i:i], &t.nodes[i]
		cp.key = n.key
		cp.value = n.value
		cp.edges = t.copyEdges(n.edges)
		return
	}
	panic("node not preallocated")
}

func (t *copyTxn) copyEdges(es edges) (cp edges) {
	n := len(es)
	i := len(t.edges) - n
	if i >= 0 {
		t.edges, cp = t.edges[:i:i], t.edges[i:]
		copy(cp, es)
		return
	}
	panic("edges not preallocated")
}
