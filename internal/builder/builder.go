package builder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var ErrNotFound = fmt.Errorf("the status was not found")
var ErrNoChanges = fmt.Errorf("no changes were detected")
var ErrNoMakefile = fmt.Errorf("no makefile was found")

type Builder struct {
	original  *git.Repository
	clone     *git.Repository
	command   []string
	makefile  string
	buildPath string
	testPath  string
}

func NewBuilder(original, clone *git.Repository, command, makefile, buildPath, testPath string) *Builder {
	return &Builder{
		original:  original,
		clone:     clone,
		command:   splitString(command),
		makefile:  makefile,
		buildPath: buildPath,
		testPath:  testPath,
	}
}

func (b *Builder) GetStatus(sha string) (string, error) {
	_, _, statusPath := b.getPaths(sha)
	contents, err := ioutil.ReadFile(statusPath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func (b *Builder) GetResults(sha string) (map[string]string, error) {
	_, filePath, _ := b.getPaths(sha)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	files := map[string]string{
		"BuildResults.md": string(contents),
	}
	return files, nil
}

func (b *Builder) getPaths(sha string) (string, string, string) {
	folder := path.Join(b.buildPath, sha[:2])
	file := path.Join(b.buildPath, sha[2:])
	status := path.Join(b.buildPath, fmt.Sprintf("%v.status", sha[2:]))
	return folder, file, status
}
func (b *Builder) Build(sha string) error {
	folderPath, filePath, statusPath := b.getPaths(sha)
	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(statusPath, []byte("pending"), 0666)
	if err != nil {
		return err
	}

	// ignoring as it will error out later when checking out if this fails
	// and no updates is an error for some reason?
	_ = b.clone.Fetch(&git.FetchOptions{})
	tree, err := b.clone.Worktree()
	if err != nil {
		return err
	}
	hash := plumbing.NewHash(sha)
	hashes := []*plumbing.Hash{
		&hash,
	}

	err = tree.Checkout(&git.CheckoutOptions{
		Hash:  *hashes[0],
		Force: true,
	})
	if err != nil {
		return err
	}

	// find files that have changed between master and this commit.
	// we use the original repo, as origin/master in the clone
	// points to the original.
	aHash, err := b.original.ResolveRevision("refs/remotes/origin/master")
	if err != nil {
		panic(err)
	}
	hashes = append(hashes, aHash)

	var commits []*object.Commit
	for _, hash := range hashes {
		commit, err := b.clone.CommitObject(*hash)
		if err != nil {
			panic(err)
		}
		commits = append(commits, commit)
	}

	res, err := commits[0].MergeBase(commits[1])
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		panic("unable to find merge base")
	}
	p, err := res[0].Patch(commits[0])

	files := []string{}
	for _, diff := range p.FilePatches() {
		to, from := diff.Files()
		if to != nil {
			files = append(files, to.Path())
		}
		if from != nil {
			files = append(files, from.Path())
		}
	}

	if len(files) == 0 {
		return ErrNoChanges
	}

	lcp := calculateLCP(files)
	lcp = filepath.Clean(lcp)

	// need to check each section of the path to find the closest one with a Makefile
	testPath := b.testPath
	for lcp != "" {
		testPath = path.Join(b.testPath, lcp)
		makefile := path.Join(testPath, "Makefile")
		if s, err := os.Stat(makefile); err == nil && !s.IsDir() {
			break
		}
		lcp = filepath.Dir(lcp)

	}
	contents := []byte("success")
	c := b.command[0]
	args := []string{}
	if len(b.command) > 1 {
		args = b.command[1:]
	}
	fmt.Printf("running %v %s\n", testPath, b.command)
	command := exec.Command(c, args...)
	command.Dir = testPath

	// tee the output to stdout.
	buf := bytes.Buffer{}
	writer := io.MultiWriter(&buf, os.Stdout)
	command.Stdout = writer
	command.Stderr = writer

	origErr := command.Run()
	out := buf.Bytes()
	if origErr != nil {
		contents = []byte("failure")
	}

	err = ioutil.WriteFile(filePath, out, 0666)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(statusPath, contents, 0666)
	if err != nil {
		return err
	}
	return origErr
}

func calculateLCP(files []string) string {
	// find the smallest common path between all changed files
	var lcp *string
	for _, f := range files {
		// variables don't work well in loops
		file := f
		if lcp == nil {
			lcp = &file
			continue
		}
		count := len(file)
		if len(*lcp) < count {
			count = len(*lcp)
		}
		i := 0
		for ; i < count; i++ {
			if (*lcp)[i] != file[i] {
				break
			}
		}
		if i == 0 {
			return ""
		}
		sub := file[:i]
		lcp = &sub
	}
	if lcp == nil {
		return ""
	}
	return *lcp
}

func splitString(s string) []string {
	lastQuote := rune(0)
	// probably some subtle bugs with null characters
	return strings.FieldsFunc(s, func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	})
}
