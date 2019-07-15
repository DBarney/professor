package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func withTmp(f func(string)) error {
	dir, err := ioutil.TempDir("/tmp", "test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	f(dir)
	return nil
}

func TestFindExistingTree(t *testing.T) {
	p := resolvePath(filepath.Join("a", "b", "c"))
	if p != "." {
		t.Fatalf("wrong path was returned %v", p)
	}
	err := withTmp(func(dir string) {
		p := resolvePath(dir)
		if p != dir {
			t.Fatalf("wrong path was returned %v", p)
		}
		p = resolvePath(filepath.Join(dir, "a", "b"))
		if p != dir {
			t.Fatalf("wrong path was returned %v", p)
		}
	})
	if err != nil {
		t.Fatalf("unable to create tmpdir %v", err)
	}
}
