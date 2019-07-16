package storage

import (
	"testing"
)

func TestCollection(t *testing.T) {
	c := Collection{"a", 1, t, 2}
	count := 0
	c.Each(func(elem interface{}) { count++ })

	if count != 4 {
		t.Fatalf("wrong number of elements were counted")
	}

	onlyInts := func(elem interface{}) bool { _, ok := elem.(int); return ok }
	e := c.First(onlyInts)
	if e == nil {
		t.Fatalf("nil is not an int")
	}
	if e.(int) != 1 {
		t.Fatalf("wrong element was returned")
	}

	e = c.Last(onlyInts)
	if e == nil {
		t.Fatalf("nil is not an int")
	}
	if e.(int) != 2 {
		t.Fatalf("wrong element was returned")
	}
}
