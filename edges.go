package trie

import "sort"

type edges []*Node

// add inserts the given node into into the target, returning the new edges and
// true if a new slice was allocated.
func (es edges) add(t *Txn, depth int, node *Node, reverse bool) (edges, bool) {
	i, old := es.get(node.key[depth], depth)
	if old != nil {
		if node, _ := mergeNodes(t, depth+1, old, node, reverse); node != old {
			return es.insert(t, i, 1, node)
		}
		return es, false
	}
	return es.insert(t, i, 0, node)
}

// cut removes n nodes starting at i from the target, returning true if a new
// slice was allocated (currently always true).
func (es edges) cut(t *Txn, i, n int) (edges, bool) {
	cp := make(edges, len(es)-n)
	copy(cp, es[:i])
	copy(cp[i:], es[i+n:])
	if debugEnabled && !cp.valid() {
		panic("invalid edges!")
	}
	return cp, true
}

func (es edges) delete(t *Txn, depth int, k []byte) (edges, bool) {
	i, old := es.get(k[depth], depth)
	if old != nil {
		if n := old.delete(t, depth+1, k); n == nil {
			return es.cut(t, i, 1)
		} else if n != old {
			return es.insert(t, i, 1, n)
		}
	}
	return es, false
}

func (es edges) deleteString(t *Txn, depth int, k string) (edges, bool) {
	i, old := es.get(k[depth], depth)
	if old != nil {
		if n := old.deleteString(t, depth+1, k); n == nil {
			return es.cut(t, i, 1)
		} else if n != old {
			return es.insert(t, i, 1, n)
		}
	}
	return es, false
}

func (es edges) get(label byte, depth int) (int, *Node) {
	i, n := es.search(label, depth)
	if i != n {
		if nd := es[i]; nd.key[depth] == label {
			return i, nd
		}
	}
	return i, nil
}

// insert replaces skip nodes with node, starting at i. It returns the new
// edges and true if a new slice was allocated (currently always true).
func (es edges) insert(t *Txn, i, skip int, node *Node) (edges, bool) {
	if debugEnabled && len(node.key) == 0 {
		panic("key too short")
	}
	cp := make(edges, len(es)+1-skip)
	copy(cp, es[:i])
	cp[i] = node
	copy(cp[i+1:], es[i+skip:])
	return cp, true
}

func (a edges) put(t *Txn, depth int, k []byte, v interface{}, b edges) (edges, bool) {
	i, old := a.get(k[depth], depth)
	if old != nil {
		if n := old.put(t, depth+1, k, v, b); n != old {
			return a.insert(t, i, 1, n)
		}
		return a, false
	}
	return a.insert(t, i, 0, t.newNode(k, v, b))
}

func (a edges) putString(t *Txn, depth int, k string, v interface{}, b edges) (edges, bool) {
	i, old := a.get(k[depth], depth)
	if old != nil {
		if n := old.putString(t, depth+1, k, v, b); n != old {
			return a.insert(t, i, 1, n)
		}
		return a, false
	}
	return a.insert(t, i, 0, t.newNode(Key(k), v, b))
}

func (es edges) search(label byte, depth int) (i, n int) {
	n = len(es)
	i = sort.Search(n, func(idx int) bool { return es[idx].key[depth] >= label })
	return
}
