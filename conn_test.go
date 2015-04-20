package main

import "testing"

func TestRead(t *testing.T) {
	i := 1
	t.Log(getHeader(i))

	i = 858993459
	t.Log(getHeader(i))

	i = 4294967295
	t.Log(getHeader(i))
}

// func TestConnTimeout(t *testing.T) {
//
// }
