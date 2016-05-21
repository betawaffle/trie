package trie

import "reflect"

type Node struct {
	key   Key
	value interface{}
	edges edges
}

func (n *Node) Get(k []byte) interface{} {
	for i := 0; n != nil; {
		end := len(n.key)

		if len(k) < end || !n.key[i:].EqualToBytes(k[i:end]) {
			break
		}
		if len(k) == end {
			return n.value
		}
		i = end

		_, n = n.edges.get(k[i], i)
	}
	return nil
}

func (n *Node) GetString(k string) interface{} {
	for i := 0; n != nil; {
		end := len(n.key)

		if len(k) < end || !n.key[i:].EqualToString(k[i:end]) {
			break
		}
		if len(k) == end {
			return n.value
		}
		i = end

		_, n = n.edges.get(k[i], i)
	}
	return nil
}

func (n *Node) Delete(k []byte) *Node {
	if n == nil {
		return nil
	}
	return n.delete(&Txn{root: n}, 0, k)
}

func (n *Node) DeleteString(k string) *Node {
	if n == nil {
		return nil
	}
	return n.deleteString(&Txn{root: n}, 0, k)
}

func (n *Node) Put(k []byte, v interface{}) *Node {
	if n == nil {
		return &Node{key: k, value: v}
	}
	return n.put(&Txn{root: n}, 0, k, v, nil)
}

func (n *Node) PutString(k string, v interface{}) *Node {
	if n == nil {
		return &Node{key: Key(k), value: v}
	}
	return n.putString(&Txn{root: n}, 0, k, v, nil)
}

func (n *Node) delete(t *Txn, depth int, k []byte) *Node {
	d, ok := n.key.commonBytesLen(k, depth)
	if ok { // not found
		return n
	}
	if d == len(k) { // exact match
		return n.deleteValue(t)
	}
	es, modified := n.edges.delete(t, d, k)
	if !modified {
		return n
	}
	if t.isMutable(n) {
		n.edges = es
		return n
	}
	return t.newNode(n.key, n.value, es)
}

func (n *Node) deleteString(t *Txn, depth int, k string) *Node {
	d, short := n.key.commonStringLen(k, depth)
	if short { // not found
		return n
	}
	if d == len(k) { // exact match
		return n.deleteValue(t)
	}
	es, modified := n.edges.deleteString(t, d, k)
	if !modified {
		return n
	}
	if t.isMutable(n) {
		n.edges = es
		return n
	}
	return t.newNode(n.key, n.value, es)
}

func (n *Node) deleteValue(t *Txn) *Node {
	switch len(n.edges) {
	case 0:
		// t.maybeFree(n)
		return nil
	case 1:
		// t.maybeFree(n)
		return n.edges[0]
	}
	if n.value != nil {
		if !t.isMutable(n) {
			return t.newNode(n.key, nil, n.edges)
		}
		n.value = nil
	}
	return n
}

// put sets the value (and merges the edges) under a given key.
func (n *Node) put(t *Txn, depth int, k []byte, v interface{}, es edges) *Node {
	d, short := n.key.commonBytesLen(k, depth)
	if short { // split
		n, _ = split(t, d, n, t.newNode(k, v, es), false)
		return n
	}
	if d == len(k) { // exact match
		return n.set(t, v, es)
	}
	es, modified := n.edges.put(t, d, k, v, es)
	if !modified {
		return n
	}
	if t.isMutable(n) {
		n.edges = es
		return n
	}
	return t.newNode(n.key, n.value, es)
}

// putString is like put, but takes a string key.
func (n *Node) putString(t *Txn, depth int, k string, v interface{}, es edges) *Node {
	d, short := n.key.commonStringLen(k, depth)
	if short { // split
		n, _ = split(t, d, n, t.newNode(Key(k), v, es), false)
		return n
	}
	if d == len(k) { // exact match
		return n.set(t, v, es)
	}
	es, modified := n.edges.putString(t, d, k, v, es)
	if !modified {
		return n
	}
	if t.isMutable(n) {
		n.edges = es
		return n
	}
	return t.newNode(n.key, n.value, es)
}

// set updates the value and merges in the provided edges.
func (n *Node) set(t *Txn, v interface{}, e edges) *Node {
	e, side := mergeEdges(t, len(n.key), n.edges, e, false)
	if side == mergeUseA {
		if reflect.DeepEqual(n.value, v) {
			debugf("set: nothing to change")
			return n
		}
		debugf("set: no edges to add, but values not equal: %#v != %#v", n.value, v)
	}
	if t.isMutable(n) {
		debugf("set: mutating")
		n.value = v
		n.edges = e
		return n
	}
	debugf("set: creating a new node")
	return t.newNode(n.key, v, e)
}
