package main

import (
	"encoding/binary"
	"testing"
)

func TestReadMsg(t *testing.T) {
	var i uint32 = 4294967295
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	t.Log(b)
}

func TestConnTimeout(t *testing.T) {

}
