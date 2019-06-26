package builder

import (
	"testing"
)

func TestLCP(t *testing.T) {
	if calculateLCP([]string{}) != "" {
		t.Fatalf("no strings should give back empty string")
	}

	if calculateLCP([]string{"asdf"}) != "asdf" {
		t.Fatalf("a single string should be given back")
	}

	if calculateLCP([]string{"a", "a"}) != "a" {
		t.Fatalf("idential strings should return one of them")
	}

	if calculateLCP([]string{"asdf", "a"}) != "a" {
		t.Fatalf("a substring should be returned")
	}

	if lcp := calculateLCP([]string{"asdf", "bcde"}); lcp != "" {
		t.Fatalf("different strings should return nothing %v", lcp)
	}

}
