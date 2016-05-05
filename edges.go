package trie

import "sort"

type edges []*Node

func (es edges) add(t *txn, node *Node) (edges, bool) {
	i, old := es.get(node.key[t.depth], t.depth)
	if old != nil {
		t.depth++

		if n := old.merge(t, node); n != old {
			return es.insert(t, i, 1, n)
		}
		return es, false
	}
	return es.insert(t, i, 0, node)
}

func (es edges) cut(t *txn, i, n int) (edges, bool) {
	cp := t.newEdges(len(es) - n)
	copy(cp, es[:i])
	copy(cp[i:], es[i+n:])
	// if !cp.valid() {
	// 	panic("invalid edges!")
	// }
	return cp, true
}

func (es edges) delete(t *txn, k []byte) (edges, bool) {
	i, old := es.get(k[t.depth], t.depth)
	if old != nil {
		t.depth++

		if n := old.delete(t, k); n == nil {
			return es.cut(t, i, 1)
		} else if n != old {
			return es.insert(t, i, 1, n)
		}
	}
	return es, false
}

func (es edges) deleteString(t *txn, k string) (edges, bool) {
	i, old := es.get(k[t.depth], t.depth)
	if old != nil {
		t.depth++

		if n := old.deleteString(t, k); n == nil {
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

func (es edges) insert(t *txn, i, skip int, node *Node) (edges, bool) {
	// if len(node.key) == 0 {
	// 	panic("key too short")
	// }
	cp := t.newEdges(len(es) + 1 - skip)
	copy(cp, es[:i])
	cp[i] = node
	copy(cp[i+1:], es[i+skip:])
	return cp, true
}

func (a edges) merge(t *txn, b edges) (edges, bool) {
	if len(b) == 0 {
		return a, true
	}
	if len(a) == 0 {
		return b, false
	}
	ns := getNodeSet(t.depth, a)
	ns.merge(t, b)
	return ns.edges(), false
}

func (a edges) put(t *txn, k []byte, v interface{}, b edges) (edges, bool) {
	i, old := a.get(k[t.depth], t.depth)
	if old != nil {
		t.depth++

		if n := old.put(t, k, v, b); n != old {
			return a.insert(t, i, 1, n)
		}
		return a, false
	}
	return a.insert(t, i, 0, t.newNode(k, v, b))
}

func (a edges) putString(t *txn, k string, v interface{}, b edges) (edges, bool) {
	i, old := a.get(k[t.depth], t.depth)
	if old != nil {
		t.depth++

		if n := old.putString(t, k, v, b); n != old {
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

func (es edges) valid() bool {
	for _, e := range es {
		if e == nil {
			return false
		}
	}
	return true
}

var nodeSets = make(chan *nodeSet, 10)

type nodeSet struct {
	nodes [256]*Node
	num   int
}

func getNodeSet(depth int, es edges) (ns *nodeSet) {
	select {
	case ns = <-nodeSets:
	default:
	}
	if ns == nil {
		ns = new(nodeSet)
	}
	ns.num = len(es)

	for _, n := range es {
		ns.nodes[n.key[depth]] = n
	}

	// panic("new nodeSet")
	return ns
}

func (ns *nodeSet) edges() edges {
	es := make(edges, 0, ns.num)
	for _, n := range ns.nodes {
		if n == nil {
			continue
		}
		es = append(es, n)
	}
	ns.nodes = [256]*Node{}
	select {
	case nodeSets <- ns:
	default:
	}
	return es
}

func (ns *nodeSet) merge(t *txn, es edges) {
	for _, n := range es {
		slot := &ns.nodes[n.key[t.depth]]

		if old := *slot; old != nil {
			if n == nil {
				ns.num--
				*slot = n
			} else {
				*slot = old.merge(t, n)
			}
		} else if n != nil {
			*slot = n
			ns.num++
		}
	}
}
