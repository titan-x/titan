package randstr

import "testing"

func TestGet(t *testing.T) {
	for i := 0; i < 50; i++ {
		l := 96
		str := Get(l)
		t.Logf("Generated string: %v", str)

		if len(str) != l {
			t.Fatalf("Expected a random string of length %v but got %v", l, len(str))
		}

		// this many collisions can't be
		if str[1] == str[2] && str[3] == str[4] && str[5] == str[6] && str[7] == str[8] && str[9] == str[10] {
			t.Fatal("Expected a random string, got repeated characters")
		}
	}
}

func TestAlphabet(t *testing.T) {
	for i := 0; i < 50; i++ {
		Alphabet = []rune("abc")
		str := Get(4)
		t.Logf("Generated string: %v", str)

		for _, s := range str {
			if string(s) != "a" && string(s) != "b" && string(s) != "c" {
				t.Fatal("Got characters outside of the given alphabet")
			}
		}
	}
}
