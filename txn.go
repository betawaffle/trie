package trie

type Txn struct {
	root *Node
	mut  map[*Node]bool
}

func (t *Txn) Prealloc(n int) {
	t.mut = make(map[*Node]bool, n)
}

func (t *Txn) Commit() *Node {
	t.mut = nil
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
	t.root, _ = mergeNodes(t, 0, t.root, n, false)
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

func (t *Txn) isMutable(n *Node) bool {
	if t == nil || t.mut == nil {
		return false
	}
	return t.mut[n]
}

func (t *Txn) newNode(k Key, v interface{}, es edges) (n *Node) {
	n = &Node{
		key:   k,
		value: v,
		edges: es,
	}
	if t != nil {
		if t.mut == nil {
			t.mut = make(map[*Node]bool)
		}
		t.mut[n] = true
	}
	return
}
