package trie

import "reflect"

func (a *Node) Merge(b *Node) (n *Node) {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	n, _ = mergeNodes(nil, 0, a, b, false)
	return
}

type mergeSide uint8

const (
	mergeNewC mergeSide = iota // a new item was returned
	mergeUseA                  // no changes were needed to the A side
	mergeUseB                  // no changes were needed to the B side
	mergeUseE                  // no changes were needed to either side
)

func mergeEdges(t *Txn, depth int, a, b edges, reverse bool) (edges, mergeSide) {
	if len(b) == 0 {
		if len(a) == 0 {
			return nil, mergeUseE
		}
		return a, mergeUseA
	}
	if len(a) == 0 {
		return b, mergeUseB
	}

	var (
		i     int
		count int
		nodes = new([256]*Node)
		side  = mergeUseE
	)

top:
	for _, m := range a {
		mLabel := m.key[depth]
		for _, n := range b[i:] {
			nLabel := n.key[depth]

			if mLabel < nLabel {
				break
			}

			if mLabel > nLabel {
				nodes[nLabel] = n
				count++
				side &= ^mergeUseA

				i++
				continue
			}

			nd, s := mergeNodes(t, depth+1, m, n, reverse)
			switch s {
			case mergeUseA:
				side &= ^mergeUseB
			case mergeUseB:
				side &= ^mergeUseA
			case mergeUseE:
				switch side {
				case mergeUseA:
					nd = m
				case mergeUseB:
					nd = n
				}
			default:
				side = mergeNewC
			}
			nodes[mLabel] = nd
			count++

			i++
			continue top
		}

		nodes[mLabel] = m
		count++
		side &= ^mergeUseB
	}

	for _, n := range b[i:] {
		nLabel := n.key[depth]

		nodes[nLabel] = n
		count++
		side &= ^mergeUseA
	}

	es := make(edges, 0, count)
	for _, n := range nodes {
		if n == nil {
			continue
		}
		es = append(es, n)
	}
	return es, side
}

func mergeNodes(t *Txn, depth int, a, b *Node, reverse bool) (*Node, mergeSide) {
	d, short := a.key.commonBytesLen(b.key, depth)
	if short {
		// one key is a prefix of the other
		// d is the prefix length
		return split(t, d, a, b, reverse)
	}

	if d < len(b.key) {
		// the keys don't match
		es, modified := a.edges.add(t, d, b, reverse)
		if !modified {
			return a, mergeUseA
		}
		if t.isMutable(a) {
			a.edges = es
			return a, mergeUseA
		}
		return t.newNode(a.key, a.value, es), mergeNewC
	}
	if debugEnabled && d > len(b.key) {
		panic("merge: sanity check failed; d > len(b.key)")
	}

	// both keys match exactly
	// time to pick the value and merge the edges

	v, vSide := mergeValues(a.value, b.value, reverse)
	e, eSide := mergeEdges(t, d, a.edges, b.edges, reverse)

	side := resolveSide(vSide, eSide)
	switch side {
	case mergeUseE:
		if reverse {
			return a, side
		}
		return b, side
	case mergeUseB:
		debugf("merge: reusing B")
		return b, side
	case mergeUseA:
		debugf("merge: reusing A")
		return a, side
	default:
		panic("merge: invalid side")
	case mergeNewC:
		// continue below
	}
	switch {
	case t.isMutable(a):
		debugf("merge: mutating A")
		a.value = v
		a.edges = e
		return a, mergeUseA // FIXME: Is this correct?
	case t.isMutable(b):
		debugf("merge: mutating B")
		b.value = v
		b.edges = e
		return b, mergeUseB // FIXME: Is this correct?
	}
	debugf("merge: creating a new node")
	return t.newNode(a.key, v, e), side
}

func mergeValues(a, b interface{}, reverse bool) (interface{}, mergeSide) {
	switch {
	case b == nil:
		if a == nil {
			return nil, mergeUseE // Both are nil, no preference.
		}
		// Always take the non-nil value.
		return a, mergeUseA
	case a == nil:
		// Always take the non-nil value.
		return b, mergeUseB
	case reflect.DeepEqual(a, b):
		// Prefer A over B if they are the same.
		// This only matters if the caller creates a new node.
		return a, mergeUseE
	case reverse:
		// If we're doing a reverse merge, prefer A.
		return a, mergeUseA
	default:
		// In a normal merge, we prefer B.
		return b, mergeUseB
	}
}

func resolveSide(value mergeSide, edges mergeSide) mergeSide {
	if value == edges {
		return value
	}
	if value == mergeUseE {
		return edges
	}
	if edges == mergeUseE {
		return value
	}
	return mergeNewC
}

// split returns a (possibly new) Node with A and/or B as edges.
func split(t *Txn, depth int, a, b *Node, reverse bool) (*Node, mergeSide) {
	if len(b.key) == depth {
		es, modified := b.edges.add(t, depth, a, !reverse)
		if !modified {
			return b, mergeUseB
		}
		if t.isMutable(a) {
			b.edges = es
			return b, mergeUseB
		}
		return t.newNode(b.key, b.value, es), mergeNewC
	}
	if b.key[depth] < a.key[depth] {
		a, b = b, a
	}
	return t.newNode(a.key[:depth], nil, edges{a, b}), mergeNewC
}
