package trie

import (
	"crypto/rand"
	"testing"

	"github.com/hashicorp/go-immutable-radix"
)

var (
	fooBar = []byte("foo bar")
	fooBaz = []byte("foo baz")
	foo    = []byte("foo")
)

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
	x := &Node{key: fooBar}
	for i := 0; i < b.N; i++ {
		x.Delete(fooBar)
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
	tree := &Node{key: fooBar}
	for i := 0; i < b.N; i++ {
		tree.Put(fooBaz, 1)
	}
}

func BenchmarkPutLong(b *testing.B) {
	keys := randomKeys(1024)
	tree := &Node{key: fooBar}
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
	x := (&Node{key: fooBar, value: 1}).Put(fooBaz, 1)
	for i := 0; i < b.N; i++ {
		x.Put(fooBaz, 1)
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
