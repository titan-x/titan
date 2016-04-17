package randstr

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Alphabet is the letters used in creating the random string.
var Alphabet = []rune(". !abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Get generates a random string sequence of given size.
func Get(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = Alphabet[rand.Intn(len(Alphabet))]
	}

	return string(b)
}
