package trie

import (
	"fmt"
	"reflect"
)

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
	return n.delete(&txn{root: n}, k)
}

func (n *Node) DeleteString(k string) *Node {
	if n == nil {
		return nil
	}
	return n.deleteString(&txn{root: n}, k)
}

func (n *Node) Histogram() map[uint8]int {
	h := make(map[uint8]int, 256)
	n.histogram(h)
	return h
}

func (n *Node) histogram(h map[uint8]int) {
	h[uint8(len(n.edges))]++
	for _, nd := range n.edges {
		nd.histogram(h)
	}
}

func (a *Node) Merge(b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return a.merge(&txn{root: a}, b)
}

func (n *Node) Put(k []byte, v interface{}) *Node {
	if n == nil {
		return &Node{key: k, value: v}
	}
	return n.put(&txn{root: n}, k, v, nil)
}

func (n *Node) PutString(k string, v interface{}) *Node {
	if n == nil {
		return &Node{key: Key(k), value: v}
	}
	return n.putString(&txn{root: n}, k, v, nil)
}

func (n *Node) String() string {
	return fmt.Sprintf("{%q:%v}", n.key, n.edges)
}

func (n *Node) Walk(fn func(*Node) bool) {
	if !fn(n) {
		return
	}
	for _, nd := range n.edges {
		nd.Walk(fn)
	}
}

func (n *Node) delete(t *txn, k []byte) *Node {
	if n.key.bytesNeedSplit(k, t) { // not found
		return n
	}
	if t.depth == len(k) { // exact match
		return n.deleteValue(t)
	}
	es, modified := n.edges.delete(t, k)
	if !modified {
		return n
	}
	return n.setEdges(t, es)
}

func (n *Node) deleteString(t *txn, k string) *Node {
	if n.key.stringNeedSplit(k, t) { // not found
		return n
	}
	if t.depth == len(k) { // exact match
		return n.deleteValue(t)
	}
	es, modified := n.edges.deleteString(t, k)
	if !modified {
		return n
	}
	return n.setEdges(t, es)
}

func (n *Node) deleteValue(t *txn) *Node {
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

// merge is like put, but takes an existing node, which can save allocations.
func (a *Node) merge(t *txn, b *Node) *Node {
	if a.key.bytesNeedSplit(b.key, t) { // split
		return a.split(t, b)
	}
	if t.depth == len(b.key) { // exact match
		return a.set(t, b.value, b.edges)
	}
	es, modified := a.edges.add(t, b)
	if !modified {
		return a
	}
	return a.setEdges(t, es)
}

// put sets the value (and merges the edges) under a given key.
func (n *Node) put(t *txn, k []byte, v interface{}, es edges) *Node {
	if n.key.bytesNeedSplit(k, t) { // split
		return n.splitNew(t, k, v, es)
	}
	if t.depth == len(k) { // exact match
		return n.set(t, v, es)
	}
	es, modified := n.edges.put(t, k, v, es)
	if !modified {
		return n
	}
	return n.setEdges(t, es)
}

// putString is like put, but takes a string key.
func (n *Node) putString(t *txn, k string, v interface{}, es edges) *Node {
	if n.key.stringNeedSplit(k, t) { // split
		return n.splitNew(t, Key(k), v, es)
	}
	if t.depth == len(k) { // exact match
		return n.set(t, v, es)
	}
	es, modified := n.edges.putString(t, k, v, es)
	if !modified {
		return n
	}
	return n.setEdges(t, es)
}

// set updates the value and merges in the provided edges.
func (n *Node) set(t *txn, v interface{}, e edges) *Node {
	e, eq := n.edges.merge(t, e)
	if eq {
		if reflect.DeepEqual(n.value, v) {
			return n
		}
	}
	if t.isMutable(n) {
		n.value = v
		n.edges = e
		return n
	}
	return t.newNode(n.key, v, e)
}

// setEdges makes a copy of the node with a different set of edges.
func (n *Node) setEdges(t *txn, es edges) *Node {
	if t.isMutable(n) {
		n.edges = es
		return n
	}
	return t.newNode(n.key, n.value, es)
}

func (a *Node) split(t *txn, b *Node) *Node {
	if b.key[t.depth] < a.key[t.depth] {
		a, b = b, a
	}
	es := t.newEdges(2)
	es[0] = a
	es[1] = b
	return t.newNode(a.key[:t.depth], nil, es)
}

func (n *Node) splitNew(t *txn, k Key, v interface{}, es edges) *Node {
	t.preallocNodes(2)
	return n.split(t, t.newNode(k, v, es))
}
