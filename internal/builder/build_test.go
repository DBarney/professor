package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

func TestGitWorktree(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "professor-test")
	if err != nil {
		t.Fatalf("unable to create test dir %v", err)
	}
	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("unable to remove test dir %v", err)
		}
	}()

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("unable to create git dir %v", err)
	}

	name := "example-git-file"
	filename := filepath.Join(dir, name)
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Fatalf("unable to create example file %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("unable to get worktree %v", err)
	}

	_, err = w.Add(name)
	if err != nil {
		t.Fatalf("unable to add file %v", err)
	}

	commit, err := w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Professor",
			Email: "prof@me.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		t.Fatalf("unable to commit changes %v", err)
	}
	err = worktreeWrap(dir, "worktree", commit.String(), func(dir string) error {
		info, err := os.Stat(dir)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("it wasn't a dir %v", info)
		}
		body, err := ioutil.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if string(body) != "hello world!" {
			return fmt.Errorf("wrong body was returend %v", string(body))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("something went wrong with the worktree %v", err)
	}
}

func TestFindMakefile(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "makefile")
	if err != nil {
		t.Fatalf("unable to make temp dir %v", err)
	}

	lcp := filepath.Join("1", "2", "3", "4")
	err = os.MkdirAll(filepath.Join(dir, lcp), 0777)
	if err != nil {
		t.Fatalf("unable to create dir tree")
	}
	path := makefileInPath(dir, lcp)
	if path != dir {
		t.Fatalf("wrong path was returned %v", path)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "1", "Makefile"), []byte("hello world!"), 0644)
	if err != nil {
		t.Fatalf("unable to create example file %v", err)
	}

	path = makefileInPath(dir, lcp)
	if path != filepath.Join(dir, "1") {
		t.Fatalf("wrong path was returned %v", path)
	}

	err = ioutil.WriteFile(filepath.Join(dir, lcp, "Makefile"), []byte("hello world!"), 0644)
	if err != nil {
		t.Fatalf("unable to create example file %v", err)
	}

	path = makefileInPath(dir, lcp)
	if path != filepath.Join(dir, lcp) {
		t.Fatalf("wrong path was returned %v", path)
	}
}
