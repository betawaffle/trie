package trie

import "testing"

func TestEdgeAdd(t *testing.T) {
	var (
		n *Node
		o = &Node{key: Key("foods"), value: 1}
		p = &Node{key: Key("foodz"), value: 2}
	)

	n = &Node{key: Key("foodie"), value: 3}

	if res, modified := (edges{o, p}).add(new(Txn), 4, n, false); !modified {
		t.Errorf("expected [%q, %q] to be modified when adding %q", o.key, p.key, n.key)
	} else {
		assertExactNode(t, n, res[0])
		assertExactNode(t, o, res[1])
		assertExactNode(t, p, res[2])
	}

	if _, modified := (edges{o, p}).add(new(Txn), 4, p, false); modified {
		t.Errorf("expected [%q, %q] to NOT be modified when adding %q", o.key, p.key, p.key)
	}

	n = &Node{key: Key("foodz"), value: 4}

	if res, modified := (edges{o, p}).add(new(Txn), 4, n, false); !modified {
		t.Errorf("expected [%q, %q] to be modified when adding %q with new value", o.key, p.key, n.key)
	} else if len(res) != 2 {
		t.Errorf("expected length to stay 2")
	} else {
		assertExactNode(t, o, res[0])
		assertExactNode(t, n, res[1])
	}
}
