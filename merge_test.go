package trie

import (
	"bytes"
	"testing"
)

// func TestMergeEdges(t *testing.T) {
// 	// testMerge := func(a, b edges) (edges, int) {
// 	// 	return merge(new(Txn), 0, a, b, false)
// 	// }
// }
//
// func TestMergeNodes_complete(t *testing.T) {
//
// }

func TestMergeNodes_nilValue(t *testing.T) {
	n := node("food", 1, edges{
		node("foodz", 2, nil),
	})
	o := node("foo", 1, edges{
		node("food", nil, edges{
			node("foods", 2, nil),
			node("foodz", 3, nil),
		}),
	})
	p := n.Merge(o)
	q := o.Merge(n)

	expectedP := node("foo", 1, edges{
		node("food", 1, edges{
			node("foods", 2, nil),
			node("foodz", 3, nil),
		}),
	})

	expectedQ := node("foo", 1, edges{
		node("food", 1, edges{
			node("foods", 2, nil),
			node("foodz", 2, nil),
		}),
	})

	if !Equal(p, expectedP) {
		t.Errorf("expected %#v, got %#v", expectedP, p)
	}

	if !Equal(q, expectedQ) {
		t.Errorf("expected %#v, got %#v", expectedQ, q)
	}
}

func TestMergeNodes(t *testing.T) {
	n := node("foodz", 4, nil)
	o := node("foodz", 2, nil)

	if res := o.Merge(n); res != n {
		t.Errorf("expected %q.merge(%q) to result in %p, got %p", o.key, n.key, n, res)
	}

	n = node("food", 1, edges{
		node("foodi", nil, edges{
			node("foodie", 2, edges{
				node("foodies", 3, nil),
			}),
			node("fooding", 4, nil),
		}),
		node("foodz", 5, nil),
	})
	o = node("foo", 1, edges{
		node("foo bar", 2, nil),
		node("food", nil, edges{
			node("fooding", 3, nil),
			node("foods", 4, edges{
				node("foods with friends", 5, nil),
			}),
		}),
	})

	p := n.Merge(o)
	q := o.Merge(n)

	expectedP := node("foo", 1, edges{
		node("foo bar", 2, nil),
		node("food", 1, edges{
			node("foodi", nil, edges{
				node("foodie", 2, edges{
					node("foodies", 3, nil),
				}),
				node("fooding", 3, nil),
			}),
			node("foods", 4, edges{
				node("foods with friends", 5, nil),
			}),
			node("foodz", 5, nil),
		}),
	})

	expectedQ := node("foo", 1, edges{
		node("foo bar", 2, nil),
		node("food", 1, edges{
			node("foodi", nil, edges{
				node("foodie", 2, edges{
					node("foodies", 3, nil),
				}),
				node("fooding", 4, nil),
			}),
			node("foods", 4, edges{
				node("foods with friends", 5, nil),
			}),
			node("foodz", 5, nil),
		}),
	})

	if !Equal(p, expectedP) {
		t.Errorf("expected %#v, got %#v", expectedP, p)
	}

	if !Equal(q, expectedQ) {
		t.Errorf("expected %#v, got %#v", expectedQ, q)
	}
}

func TestMerge_1(t *testing.T) {
	n := &Node{key: fooBar, value: []byte("12345")}
	m := n.Put(fooBaz, []byte("2"))

	tx := new(Txn)
	tx.Put(fooBar, []byte("12345"))
	tx.Put(fooBaz, []byte("2"))
	tx.Merge(m)

	o := tx.Commit()

	if res := o.GetString("foo bar").([]byte); !bytes.Equal(res, []byte("12345")) {
		t.Errorf(`expected "12345", got %q`, res)
	}
	if res := o.GetString("foo baz").([]byte); !bytes.Equal(res, []byte("2")) {
		t.Errorf(`expected "2", got %q`, res)
	}
}

func TestMerge_2(t *testing.T) {
	n := node("foo bar", []byte("12345"), nil)
	o := n.Merge(n.PutString("foo baz", []byte("2")))

	if res := o.GetString("foo bar").([]byte); !bytes.Equal(res, []byte("12345")) {
		t.Fatalf(`expected "12345", got %q`, res)
	}
	if res := o.GetString("foo baz").([]byte); !bytes.Equal(res, []byte("2")) {
		t.Fatalf(`expected "2", got %q`, res)
	}
}

func TestMerge_3(t *testing.T) {
	n := node("foo bar", []byte("12345"), nil)
	o := n.Merge(n.PutString("foo", []byte("2")))

	if res := o.GetString("foo bar").([]byte); !bytes.Equal(res, []byte("12345")) {
		t.Fatalf(`expected "12345", got %q`, res)
	}
	if res := o.GetString("foo").([]byte); !bytes.Equal(res, []byte("2")) {
		t.Fatalf(`expected "2", got %q`, res)
	}
}
