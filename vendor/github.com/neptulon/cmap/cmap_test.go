package cmap

import "testing"

func TestMap(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("newly created map is nil")
	}

	if m.Len() != 0 {
		t.Fatal("newly created map is not empty")
	}

	m.Set("foo", "bar")

	if val, ok := m.GetOk("foo"); ok {
		if bar := val.(string); bar != "bar" {
			t.Fatal("retrieved value is not 'bar'")
		}
	} else {
		t.Fatal("cannot value using key 'foo' in map")
	}

	m.Delete("foo")

	if m.Len() != 0 {
		t.Fatal("failed to delte 'foo' from map")
	}
}

func TestRange(t *testing.T) {
	m := New()
	m.Set("foo", "bar")
	m.Set("foo2", "bar2")

	var i int
	m.Range(func(val interface{}) {
		i++
		if val != "bar" && val != "bar2" {
			t.Fatal("invalid value retrieved during iteration")
		}
	})

	if i != 2 {
		t.Fatal("invalid number of elements were iterated")
	}
}
