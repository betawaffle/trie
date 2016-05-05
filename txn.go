package trie

import "fmt"

type txn struct {
	root  *Node
	depth int
	mut   map[*Node]bool
	nodes []Node
	edges []*Node

	// stats
	newNodes int
}

func (t *txn) PrintHistogram() {
	// fmt.Printf("mutable nodes: %d, new nodes: %d\n", len(t.mut), t.newNodes)

	h := t.root.Histogram()
	for i := 0; i < 256; i++ {
		n := h[uint8(i)]
		if n == 0 {
			continue
		}
		fmt.Printf("%3d: %d\n", i, n)
	}
}

func newTxn(n int) *txn {
	return &txn{
		mut:   make(map[*Node]bool, n),
		nodes: make([]Node, n),
		edges: make([]*Node, n),
	}
}

func (t *txn) Prealloc(n int) {
	t.mut = make(map[*Node]bool, n)
	t.nodes = make([]Node, n)
	t.edges = make([]*Node, n-1)
}

func (t *txn) Commit() *Node {
	t.mut = nil
	t.nodes = nil
	t.edges = nil
	return t.root
}

func (t *txn) Delete(k []byte) {
	if t.root != nil {
		t.depth = 0
		t.root = t.root.delete(t, k)
	}
}

func (t *txn) DeleteString(k string) {
	if t.root != nil {
		t.depth = 0
		t.root = t.root.deleteString(t, k)
	}
}

func (t *txn) Merge(n *Node) {
	if n == nil {
		return
	}
	if t.root == nil {
		t.root = n
		return
	}
	t.depth = 0
	t.root = t.root.merge(t, n)
}

func (t *txn) Put(k []byte, v interface{}) {
	if t.root != nil {
		t.depth = 0
		t.root = t.root.put(t, k, v, nil)
		return
	}
	t.root = t.newNode(k, v, nil)
}

func (t *txn) PutString(k string, v interface{}) {
	if t.root != nil {
		t.depth = 0
		t.root = t.root.putString(t, k, v, nil)
		return
	}
	t.root = t.newNode(Key(k), v, nil)
}

func (t *txn) preallocNodes(n int) {
	if len(t.nodes) < n {
		t.nodes = make([]Node, n)
	}
}

func (t *txn) isMutable(n *Node) bool {
	if t.mut == nil {
		return false
	}
	return t.mut[n]
}

func (t *txn) newNode(k Key, v interface{}, es edges) (n *Node) {
	// if t.nodes == nil {
	// 	t.nodes = make([]Node, 8)
	// }
	if i := len(t.nodes) - 1; i >= 0 {
		t.nodes, n = t.nodes[:i:i], &t.nodes[i]
		n.key = k
		n.value = v
		n.edges = es
	} else {
		// panic("node not preallocated")
		n = &Node{
			key:   k,
			value: v,
			edges: es,
		}
		t.newNodes++
	}
	if t.mut == nil {
		t.mut = make(map[*Node]bool)
	}
	t.mut[n] = true
	return
}

func (t *txn) newEdges(n int) (es edges) {
	if len(t.edges) < 10 {
		t.edges = make([]*Node, 256)
	}
	if i := len(t.edges) - n; i >= 0 {
		t.edges, es = t.edges[:i:i], t.edges[i:]
		return
	}
	// println(fmt.Sprintf("allocating %d edges", n))
	return make(edges, n)
}
