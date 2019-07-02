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
	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var ErrNotFound = fmt.Errorf("the status was not found")
var ErrNoChanges = fmt.Errorf("no changes were detected")
var ErrNoMakefile = fmt.Errorf("no makefile was found")

var targetMatch = regexp.MustCompile(" +(.+) `([^']+)'")

type Type int32

type filter struct {
	w io.Writer
}

const (
	Start Type = iota
	Stop
	Message
	Target
)

type BuildEvent struct {
	T       Type
	Message string
}

type Builder struct {
	original  *git.Repository
	makefile  string
	target    string
	buildPath string
	testPath  string

	publish chan *BuildEvent
}

func NewBuilder(original *git.Repository, makefile, target, buildPath, testPath string) *Builder {
	return &Builder{
		original:  original,
		makefile:  makefile,
		target:    target,
		buildPath: buildPath,
		testPath:  testPath,
	}
}

func (b *Builder) GetStatus(sha string) (string, error) {
	_, _, statusPath, _ := b.getPaths(sha)
	contents, err := ioutil.ReadFile(statusPath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func (b *Builder) GetResults(sha string) (map[string]string, error) {
	_, filePath, _, _ := b.getPaths(sha)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	files := map[string]string{
		"BuildResults.md": string(contents),
	}
	return files, nil
}

func (b *Builder) GetStream() <-chan *BuildEvent {
	if b.publish == nil {
		b.publish = make(chan *BuildEvent)
	}
	return b.publish
}

func (b *Builder) getPaths(sha string) (string, string, string, string) {
	folder := path.Join(b.buildPath, sha[:2])
	file := path.Join(b.buildPath, sha[:2], sha[2:])
	status := path.Join(b.buildPath, fmt.Sprintf("%v.status", sha[2:]))
	worktree := path.Join(b.testPath, sha[:2], sha[2:])
	return folder, file, status, worktree
}

func (b *Builder) update(t Type, m string) {
	if b.publish == nil {
		return
	}
	b.publish <- &BuildEvent{
		T:       t,
		Message: m,
	}
}
func (b *Builder) Build(sha string) error {
	b.update(Start, sha)
	defer b.update(Stop, sha)
	folderPath, filePath, statusPath, worktree := b.getPaths(sha)
	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(statusPath, []byte("pending"), 0666)
	if err != nil {
		return err
	}

	// create a temporary worktree to use for the build
	// TODO: change to tmpdir
	// TODO: make PR to go-git to add actual worktree support
	o, err := exec.Command("git", "worktree", "add", worktree, sha).CombinedOutput()
	fmt.Printf("create tree: %v %v\n", string(o), err)
	if err != nil {
		return err
	}
	defer func() {
		e := os.RemoveAll(worktree)

		o, err := exec.Command("git", "worktree", "prune").CombinedOutput()
		fmt.Printf("create tree: %v %v %v\n", string(o), e, err)
		fmt.Printf("cleanup: %v\n", e)
	}()

	hash := plumbing.NewHash(sha)
	hashes := []*plumbing.Hash{
		&hash,
	}
	// find files that have changed between master and this commit.
	aHash, err := b.original.ResolveRevision("refs/remotes/origin/master")
	if err != nil {
		panic(err)
	}
	hashes = append(hashes, aHash)

	var commits []*object.Commit
	for _, hash := range hashes {
		commit, err := b.original.CommitObject(*hash)
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
	testPath := fmt.Sprintf("%v%v", b.testPath, b.makefile)
	// if there isn't a specific makfile, search for one
	if b.makefile == "" {
		lcp := calculateLCP(files)
		lcp = filepath.Clean(lcp)
		fmt.Printf("searching for a Makefile, starting at %v\n", lcp)

		// need to check each section of the path to find the closest one with a Makefile
		for lcp != "" {
			testPath = path.Join(worktree, lcp)
			makefile := path.Join(testPath, "Makefile")
			if s, err := os.Stat(makefile); err == nil && !s.IsDir() {
				break
			}
			lcp = filepath.Dir(lcp)

		}
	}
	contents := []byte("success")
	fmt.Printf("running %v make %s\n", testPath, b.target)
	command := exec.Command("make", b.target, "--debug=b")
	command.Dir = testPath

	// tee the output to stdout and to the publish channel
	buf := bytes.Buffer{}
	writer := io.MultiWriter(&buf, os.Stdout)
	writer = io.MultiWriter(filter{w: writer}, b)
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

func (f filter) Write(p []byte) (int, error) {
	lines := bytes.Split(p, []byte("\n"))
	skipped := 0
	for i := 0; i < len(lines); {
		l := lines[i]
		match := targetMatch.FindSubmatch(l)
		if len(match) == 0 {
			i++
			continue
		}
		skipped += len(l) + 1 // +1 because of the new line character
		lines = append(lines[:i], lines[i+1:]...)
	}
	p1 := bytes.Join(lines, []byte("\n"))
	n, err := f.w.Write(p1)
	return n + skipped, err
}

func (b *Builder) Write(p []byte) (int, error) {
	lines := bytes.Split(p, []byte("\n"))
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		match := targetMatch.FindSubmatch(l)
		if len(match) != 0 {
			b.update(Target, string(match[2]))
			continue
		}
		b.update(Message, string(l))
	}
	return len(p), nil
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
