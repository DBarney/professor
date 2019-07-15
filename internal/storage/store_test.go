package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

type Info struct {
	Field string
}

type Status struct {
	Code int
}

func withDir(f func(string)) error {
	dir, err := ioutil.TempDir("/tmp", "store")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	f(dir)
	return nil
}

func TestStore(t *testing.T) {
	err := withDir(func(dir string) {
		s, err := New(dir)
		if err != nil {
			t.Fatalf("unable to open store %v", err)
		}

		err = s.Create("key", Info{Field: "info field\n"})
		if err != nil {
			t.Fatalf("unable to create entry %v", err)
		}

		err = s.Create("key", Info{Field: "info field\n"})
		if err != ErrAlreadyExists {
			t.Fatalf("should not overwrite entries %v", err)
		}

		err = s.Append("key", Status{Code: 99})
		if err != nil {
			t.Fatalf("unable to append information %v", err)
		}

		entries, err := s.Get("key")
		if err != nil {
			t.Fatalf("unable to get entries %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("wrong number of entries was were returned %v", entries)
		}
	})
	if err != nil {
		t.Fatalf("unable to setup dir %v", err)
	}
}
