package trie

import (
	"crypto/rand"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-immutable-radix"
)

var (
	k1 = []byte("foo bar")
	k2 = []byte("foo baz")
	k3 = []byte("foo")
)

func TestSizes(t *testing.T) {
	if size := unsafe.Sizeof(Node{}); size != 64 {
		t.Errorf("expected node size to be 64, got %d", size)
	}
}

func TestPut(t *testing.T) {
	n := &Node{key: k1, value: []byte("12345")}
	m := n.Put(k2, []byte("2"))

	if !m.key.EqualToString("foo ba") {
		t.Fatalf(`expected "foo ba", got %q`, m.key)
	}
	if len(m.edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(m.edges))
	}
	if e := m.edges[0]; e != n {
		t.Errorf("expected first node to be reused")
	}
}

// func ExampleReport() {
// 	keys := randomKeys(1024)
//
// 	tx := newTxn(1024)
// 	for _, k := range keys {
// 		tx.Put(k, 1)
// 	}
// 	tx.PrintHistogram()
// 	//   0: 1024
// 	//   2: 46
// 	//   3: 55
// 	//   4: 50
// 	//   5: 37
// 	//   6: 26
// 	//   7: 17
// 	//   8: 4
// 	//   9: 5
// 	//  10: 1
// 	//  12: 1
// 	// 250: 1
// }

// func TestDelete(t *testing.T) {
// 	keys := randomKeys(1024)
// 	tree := buildTree(keys)
//
// 	nodes, edges := countTree(tree)
// 	fmt.Printf("nodes: %d, edges: %d\n", nodes, edges)
//
// 	_ = t
// }

func BenchmarkGetLong(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys[i%1024])
	}
}

func BenchmarkGetLongCompetitor(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildCompetitorTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys[i%1024])
	}
}

func BenchmarkDelete(b *testing.B) {
	x := &Node{key: k1}
	for i := 0; i < b.N; i++ {
		x.Delete(k1)
	}
}

func BenchmarkDeleteLong(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Delete(keys[i%1024])
	}
}

func BenchmarkDeleteLongCompetitor(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildCompetitorTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Delete(keys[i%1024])
	}
}

func BenchmarkDenseCopy(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.DenseCopy()
	}
}

func BenchmarkPutNew(b *testing.B) {
	tree := &Node{key: k1}
	for i := 0; i < b.N; i++ {
		tree.Put(k2, 1)
	}
}

func BenchmarkPutLong(b *testing.B) {
	keys := randomKeys(1024)
	tree := &Node{key: k1}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree = tree.Put(keys[i%1024], 1)
	}
}

func BenchmarkPutLongCompetitor(b *testing.B) {
	keys := randomKeys(1024)
	tree := iradix.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree, _, _ = tree.Insert(keys[i%1024], 1)
	}
}

func BenchmarkPutExisting(b *testing.B) {
	x := &Node{key: k1}
	y := x.Put(k2, 1)
	for i := 0; i < b.N; i++ {
		y.Put(k2, 1)
	}
}

func BenchmarkPutExistingLong(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Put(keys[i%1024], 1)
	}
}

func BenchmarkPutExistingLongCompetitor(b *testing.B) {
	keys := randomKeys(1024)
	tree := buildCompetitorTree(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i%1024], 1)
	}
}

func BenchmarkPutExistingLongMap(b *testing.B) {
	keys := randomKeys(1024)
	m := buildMap(keys)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m[string(keys[i%1024])] = 1
	}
}

func BenchmarkRandomTree(b *testing.B) {
	keys := randomKeys(1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buildTree(keys)
	}
}

func BenchmarkRandomCompetitorTree(b *testing.B) {
	keys := randomKeys(1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buildCompetitorTree(keys)
	}
}

func BenchmarkRandomMap(b *testing.B) {
	keys := randomKeys(1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buildMap(keys)
	}
}

func countTree(root *Node) (nodes, edges int) {
	root.Walk(func(n *Node) bool {
		nodes += 1
		edges += len(n.edges)
		return true
	})
	return
}

func buildMap(keys [][]byte) map[string]int {
	m := make(map[string]int)
	for _, k := range keys {
		m[string(k)] = 1
	}
	return m
}

func buildTree(keys [][]byte) *Node {
	tx := new(Txn)
	tx.Prealloc(len(keys) + 300)
	// newTxn(len(keys) + 300)
	for _, k := range keys {
		tx.Put(k, 1)
	}
	return tx.Commit()
}

func buildCompetitorTree(keys [][]byte) *iradix.Tree {
	tx := iradix.New().Txn()
	for _, k := range keys {
		tx.Insert(k, 1)
	}
	return tx.Commit()
}

func randomKeys(n int) [][]byte {
	data := make([]byte, 8192)
	rand.Read(data)

	keys := make([][]byte, n)
	for i := 0; i < n; i++ {
		var (
			start = (i * 2) % 8000
			end   = (start + 256)
		)
		if end > len(data) {
			end = len(data)
		}
		keys[i] = data[start:end]
	}
	return keys
}
