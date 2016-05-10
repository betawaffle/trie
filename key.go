package trie

import "bytes"

type Key []byte

func (k Key) CommonBytesLen(b []byte) (int, bool) {
	i, short := k.commonBytesLen(b, 0)
	return i, !short
}

func (k Key) CommonStringLen(s string) (int, bool) {
	i, short := k.commonStringLen(s, 0)
	return i, !short
}

func (k Key) EqualToBytes(b []byte) bool {
	return bytes.Equal(k, b)
}

func (k Key) EqualToString(s string) bool {
	return string(k) == s // cool, this doesn't copy!
}

func (k Key) commonBytesLen(b []byte, depth int) (n int, short bool) {
	if len(b) < len(k) {
		k, short = k[:len(b)], true
	}
	for i, c := range k[depth:] {
		i += depth

		if b[i] != c {
			return i, true
		}
	}
	return len(k), short
}

func (k Key) commonStringLen(s string, depth int) (n int, short bool) {
	if len(s) < len(k) {
		k, short = k[:len(s)], true
	}
	for i, c := range k[depth:] {
		i += depth

		if s[i] != c {
			return i, true
		}
	}
	return len(k), short
}

func (k Key) bytesNeedSplit(b []byte, t *Txn) (short bool) {
	if len(b) < len(k) {
		k, short = k[:len(b)], true
	}
	b = b[t.depth:]

	for i, c := range k[t.depth:] {
		if b[i] != c {
			t.depth += i
			return true
		}
	}
	t.depth = len(k)
	return short
}

func (k Key) stringNeedSplit(s string, t *Txn) (short bool) {
	if len(s) < len(k) {
		k, short = k[:len(s)], true
	}
	s = s[t.depth:]

	for i, c := range k[t.depth:] {
		if s[i] != c {
			t.depth += i
			return true
		}
	}
	t.depth = len(k)
	return short
}
