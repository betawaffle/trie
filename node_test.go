package trie

import (
	"bytes"
	"runtime"
	"testing"
	"unsafe"
)

var emptyTrie *Node

var foodTrie = node("food", nil, edges{
	node("foodie", 2, edges{
		node("foodies", 3, nil),
	}),
	node("foods", 4, nil),
	node("foodz", 5, nil),
})

var footTrie = node("foot", nil, edges{
	node("footie", 2, edges{
		node("footies", 3, nil),
	}),
	node("foots", 4, nil),
	node("footz", 5, nil),
})

func TestSizeOfNode(t *testing.T) {
	if size := unsafe.Sizeof(Node{}); size != 64 {
		t.Errorf("expected Node to be 64 bytes, got %d", size)
	}
}

func TestNodeGet(t *testing.T) {
	testGet := func(n *Node, key string) interface{} {
		a := n.Get([]byte(key))
		b := n.GetString(key)

		if a != b {
			t.Errorf("Get does not agree with GetString")
		}

		return a
	}

	if testGet(foodTrie, "foo") != nil {
		t.Errorf(`"foo" should not have been found`)
	}

	if testGet(foodTrie, "foods") != 4 {
		t.Errorf(`"foods" should have been 4`)
	}
}

func TestNodePut(t *testing.T) {
	testPut := func(n *Node, key string, val interface{}) *Node {
		a := n.Put([]byte(key), val)
		b := n.PutString(key, val)

		if !Equal(a, b) {
			t.Errorf("Put does not agree with PutString for %q with trie: %#v", key, n)
		}

		return a
	}

	if testPut(emptyTrie, "foodie", 2) == nil {
		t.Errorf("expected Put to return a new Node")
	}

	if testPut(foodTrie, "foodie", 2) != foodTrie {
		t.Errorf("expected Put to do nothing")
	}

	if testPut(foodTrie, "foodies", 2) == foodTrie {
		t.Errorf("expected Put to replace a Node")
	}

	if testPut(footTrie, "food", 2) == footTrie {
		t.Errorf("expected Put to replace a Node")
	}

	if testPut(foodTrie, "foodly", 2) == foodTrie {
		t.Errorf("expected Put to return a new Node")
	}

	n := testPut(emptyTrie, "foo bar", []byte("12345"))
	m := testPut(n, "foo baz", []byte("2"))

	if !m.key.EqualToString("foo ba") {
		t.Errorf(`expected "foo ba", got %q`, m.key)
	}
	if len(m.edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(m.edges))
	}
	if e := m.edges[0]; e != n {
		t.Errorf("expected first node to be reused")
	}
}

func TestNodeDelete(t *testing.T) {
	testDelete := func(n *Node, key string) *Node {
		a := n.Delete([]byte(key))
		b := n.DeleteString(key)

		if !Equal(a, b) {
			t.Errorf("Delete does not agree with DeleteString for %q with trie: %#v", key, n)
		}

		return a
	}

	if testDelete(emptyTrie, "food") != nil {
		t.Errorf("expected Delete to return nil")
	}

	if testDelete(foodTrie, "foo") != foodTrie {
		t.Errorf("expected Delete to do nothing")
	}

	if testDelete(foodTrie, "foodsies") != foodTrie {
		t.Errorf("expected Delete to do nothing")
	}

	if testDelete(foodTrie, "foodie") == foodTrie {
		t.Errorf("expected Delete to return a new Node")
	}

	if testDelete(foodTrie, "foodies") == foodTrie {
		t.Errorf("expected Delete to return a new Node")
	}

	n := node("foo bar", []byte("12345"), nil)
	m := n.PutString("foo baz", []byte("2")).DeleteString("foo bar")

	if res := m.GetString("foo bar"); res != nil {
		t.Errorf(`expected nil, got %q`, res)
	}
	if res := m.GetString("foo baz").([]byte); !bytes.Equal(res, []byte("2")) {
		t.Errorf(`expected "2", got %q`, res)
	}
}

func assertExactNode(t *testing.T, expected, actual *Node) {
	if actual != expected {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("expected %q (%p), got %q (%p) at %s:%d", expected.key, expected, actual.key, actual, file, line)
	}
}

func node(key string, val interface{}, es edges) *Node {
	return &Node{
		key:   Key(key),
		value: val,
		edges: es,
	}
}
