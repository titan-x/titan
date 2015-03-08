package main

import "testing"

func TestRandString(t *testing.T) {
	l := 12304
	str := randString(l)

	if len(str) != l {
		t.Fatalf("Expected a random string of length %v but got %v", l, len(str))
	}
	if str[1] == str[2] && str[3] == str[4] && str[5] == str[6] && str[7] == str[8] {
		t.Fatal("Expected a random string, got repeated characters.")
	}
}

func TestGetID(t *testing.T) {
	for i := 0; i < 50; i++ {
		id, err := getID()
		if err != nil {
			t.Fatalf("Error while generating unique ID: %v", err)
		}
		if len(id) != 26 {
			t.Fatalf("Expected a string of length 26 but got %v", len(id))
		}
		if id[3] == id[4] && id[5] == id[6] && id[7] == id[8] && id[9] == id[10] {
			t.Fatal("Expected a random string, got repeated characters.")
		}
	}
}
