package trie

import (
	"fmt"
	"strings"
)

func (es edges) GoString() string {
	if es == nil {
		return "edges(nil)"
	}
	if len(es) == 0 {
		return "edges{}"
	}
	strs := make([]string, 0, len(es))
	for _, n := range es {
		ns := strings.Replace(n.GoString(), "\n", "\n\t", -1)
		strs = append(strs, ns)
	}
	return "edges{\n\t" + strings.Join(strs, ",\n\t") + ",\n}"
}

func (es edges) valid() bool {
	for _, e := range es {
		if e == nil {
			return false
		}
	}
	return true
}

func (n *Node) GoString() string {
	if n == nil {
		return "(*Node)(nil)"
	}
	es := strings.Replace(n.edges.GoString(), "\n", "\n\t", -1)
	if n.value == nil {
		return fmt.Sprintf("&Node{\n\tkey:   Key(%q),\n\tedges: %s,\n}", n.key, es)
	}
	return fmt.Sprintf("&Node{\n\tkey:   Key(%q),\n\tvalue: %#v,\n\tedges: %s,\n}", n.key, n.value, es)
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

func (t *Txn) PrintHistogram() {
	// fmt.Printf("mutable nodes: %d, new nodes: %d\n", len(t.mut), t.newNodes)

	h := t.root.Histogram()
	for i := 0; i < 256; i++ {
		n := h[uint8(i)]
		if n == 0 {
			continue
		}
		fmt.Printf("%3d: %d\n", i, n)
	}
}
