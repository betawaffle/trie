package trie

import "fmt"

type Txn struct {
	root  *Node
	mut   map[*Node]bool
	nodes []Node
	edges []*Node

	// stats
	newNodes int
}

func (t *Txn) PrintHistogram() {
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

func (t *Txn) Prealloc(n int) {
	t.mut = make(map[*Node]bool, n)
	t.nodes = make([]Node, n)
	t.edges = make([]*Node, n-1)
}

func (t *Txn) Commit() *Node {
	t.mut = nil
	t.nodes = nil
	t.edges = nil
	return t.root
}

func (t *Txn) Delete(k []byte) {
	if t.root != nil {
		t.root = t.root.delete(t, 0, k)
	}
}

func (t *Txn) DeleteString(k string) {
	if t.root != nil {
		t.root = t.root.deleteString(t, 0, k)
	}
}

func (t *Txn) Merge(n *Node) {
	if n == nil {
		return
	}
	if t.root == nil {
		t.root = n
		return
	}
	t.root = t.root.merge(t, 0, n)
}

func (t *Txn) Put(k []byte, v interface{}) {
	if t.root != nil {
		t.root = t.root.put(t, 0, k, v, nil)
		return
	}
	t.root = t.newNode(k, v, nil)
}

func (t *Txn) PutString(k string, v interface{}) {
	if t.root != nil {
		t.root = t.root.putString(t, 0, k, v, nil)
		return
	}
	t.root = t.newNode(Key(k), v, nil)
}

func (t *Txn) preallocNodes(n int) {
	if len(t.nodes) < n {
		t.nodes = make([]Node, n)
	}
}

func (t *Txn) isMutable(n *Node) bool {
	if t.mut == nil {
		return false
	}
	return t.mut[n]
}

func (t *Txn) newNode(k Key, v interface{}, es edges) (n *Node) {
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

func (t *Txn) newEdges(n int) (es edges) {
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
