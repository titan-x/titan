package neptulon

import "testing"

func TestRead(t *testing.T) {
	i := 1
	t.Log(makeHeaderBytes(i, 4))

	i = 858993459
	t.Log(makeHeaderBytes(i, 4))

	i = 4294967295
	t.Log(makeHeaderBytes(i, 4))
}

// func TestConnTimeout(t *testing.T) {
//
// }
