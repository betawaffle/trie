package trie

import (
	"bytes"
	"reflect"
)

func Equal(a, b *Node) bool {
	if a == b {
		return true
	}
	if !bytes.Equal(a.key, b.key) {
		debugf("Equal: different key: %q != %q", a.key, b.key)
		return false
	}
	if !reflect.DeepEqual(a.value, b.value) {
		debugf("Equal: different values for %q: %#v != %#v", a.key, a.value, b.value)
		return false
	}
	return edgesEqual(a.edges, b.edges)
}

func edgesEqual(a, b edges) bool {
	if len(a) != len(b) {
		debugf("edgesEqual: different lengths; %d != %d", len(a), len(b))
		return false
	}
	switch {
	case a == nil:
		if b == nil {
			return true
		}
		debugf("edgesEqual: nil A")
		return false
	case b == nil:
		debugf("edgesEqual: nil B")
		return false
	}
	for i, e := range a {
		if !Equal(e, b[i]) {
			return false
		}
	}
	return true
}
