// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 13, 2020

package gocache

// ByteView :holds an immutable view of bytes.
type ByteView struct {
	b []byte
}

// Len :return length of b
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice :Returns a copy of a byte slice
func (v ByteView) ByteSlice() []byte {
	c := make([]byte, len(v.b))
	copy(c, v.b)
	return c
}

// String +!
func (v ByteView) String() string {
	return string(v.b)
}
