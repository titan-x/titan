package test

import (
	"fmt"
	"testing"
)

func BenchmarkAuth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func BenchmarkClientCertAuth(b *testing.B) {
	// for various certificate key sizes (512....4096) and ECDSA, and with/without resumed handshake / session tickets
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func BenchmarkQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func BenchmarkParallelThroughput(b *testing.B) {
	// for various conn levels vs. message per second: 50:xxxx, 500:xxx, 5000:xx, ... conn/mps (hopefully!)
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
