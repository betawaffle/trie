package trie

func (n *Node) Key() []byte {
	return n.key
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) Walk(fn func(*Node) bool) {
	if n == nil {
		return
	}
	n.walk(fn)
}

func (n *Node) WalkChan(ch chan<- *Node) {
	if n != nil {
		n.walkChan(ch)
	}
	close(ch)
}

func (n *Node) walk(fn func(*Node) bool) {
	if n.value != nil && !fn(n) {
		return
	}
	for _, nd := range n.edges {
		nd.Walk(fn)
	}
}

func (n *Node) walkChan(ch chan<- *Node) {
	if n.value != nil {
		ch <- n
	}
	for _, nd := range n.edges {
		nd.walkChan(ch)
	}
}
