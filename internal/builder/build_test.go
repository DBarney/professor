package builder

import (
	"bytes"
	"testing"
)

func TestRegex(t *testing.T) {
	m := targetMatch.FindStringSubmatch(" Must remake target `test'.")
	if len(m) != 3 {
		t.Fatalf("wrong number of matches %v", m)
	}
	if m[1] != "Must remake target" {
		t.Fatalf("message wasn't trimmed correctly '%v'", m[1])
	}
}

func TestWriteFilterer(t *testing.T) {
	buf := &bytes.Buffer{}
	f := filter{w: buf}
	n, err := f.Write([]byte("asdf"))
	if n != 4 || err != nil {
		t.Fatalf("plan write should not be filtered %v, %v", n, err)
	}
	b := buf.Bytes()
	if !bytes.Equal(b, []byte("asdf")) {
		t.Fatalf("the wrong thing was written")
	}
	buf.Reset()

	n, err = f.Write([]byte("asdf\n"))
	if n != 5 || err != nil {
		t.Fatalf("single line write should not be filtered %v, %v", n, err)
	}
	b = buf.Bytes()
	if !bytes.Equal(b, []byte("asdf\n")) {
		t.Fatalf("the wrong thing was written")
	}
	buf.Reset()

	n, err = f.Write([]byte("asdf\nqwer"))
	if n != 9 || err != nil {
		t.Fatalf("multi line write should not be filtered %v, %v", n, err)
	}
	b = buf.Bytes()
	if !bytes.Equal(b, []byte("asdf\nqwer")) {
		t.Fatalf("the wrong thing was written")
	}
	buf.Reset()

	n, err = f.Write([]byte("Must remake target `test'.\nasdf"))
	if n != 31 || err != nil {
		t.Fatalf("filtered line should be removed %v %v", n, err)
	}
	b = buf.Bytes()
	if !bytes.Equal(b, []byte("asdf")) {
		t.Fatalf("the wrong thing was written")
	}
	buf.Reset()

	s := "r\n	   File `s' does not exist.\ns\n	   File `t' does not exist.\nt\n	  Must remake target `u'.\nu\n	  Must remake target `v'.\nv\n	  Must remake target `w'."
	n, err = f.Write([]byte(s))
	if n != 148 || err != nil {
		t.Fatalf("filtered line should be removed %v %v", n, err)
	}
	b = buf.Bytes()
	if !bytes.Equal(b, []byte("r\ns\nt\nu\nv")) {
		t.Fatalf("the wrong thing was written '%v'", string(b))
	}
	buf.Reset()
}

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
