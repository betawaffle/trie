package trie

import (
	"crypto/rand"
	"testing"
)

func TestCommonStringLen(t *testing.T) {
	table := []struct {
		K Key
		S string
		N int
	}{
		{Key("foo"), "bar", 0},
		{Key("bar"), "baz", 2},
		{Key("foobar"), "foo", 3},
		{Key("foo"), "foobar", 3},
	}
	for _, x := range table {
		if n, _ := x.K.CommonStringLen(x.S); n != x.N {
			t.Errorf("expected %d; got %d", x.N, n)
		}
	}
}

func TestCommonBytesLen(t *testing.T) {
	table := []struct {
		K Key
		B []byte
		N int
	}{
		{Key("foo"), []byte("bar"), 0},
		{Key("bar"), []byte("baz"), 2},
		{Key("foobar"), []byte("foo"), 3},
		{Key("foo"), []byte("foobar"), 3},
	}
	for _, x := range table {
		if n, _ := x.K.CommonBytesLen(x.B); n != x.N {
			t.Errorf("expected %d; got %d", x.N, n)
		}
	}
}

func BenchmarkEqualToBytes(b *testing.B) {
	k := make(Key, 256)
	rand.Read(k)
	x := make([]byte, len(k))
	copy(x, k)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k.EqualToBytes(x)
	}
}

func BenchmarkEqualToString(b *testing.B) {
	k := make(Key, 256)
	rand.Read(k)
	x := string(k)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k.EqualToString(x)
	}
}
