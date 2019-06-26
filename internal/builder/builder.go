package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var ErrNotFound = fmt.Errorf("the status was not found")
var ErrNoMakefile = fmt.Errorf("no makefile was found")

type Builder struct {
	original  *git.Repository
	clone     *git.Repository
	makefile  string
	buildPath string
	testPath  string
}

func NewBuilder(original, clone *git.Repository, makefile, buildPath, testPath string) *Builder {
	return &Builder{
		original:  original,
		clone:     clone,
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
		"BuildResults.txt": string(contents),
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

	fmt.Printf("%v\n", commits)

	res, err := commits[0].MergeBase(commits[1])
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		panic("unable to find merge base")
	}
	fmt.Printf("merge base: %v\n", res)
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

	fmt.Printf("files: %v\n", files)

	lcp := calculateLCP(files)
	lcp = filepath.Clean(lcp)

	fmt.Printf("got lcp: '%v'\n", lcp)
	// need to check each section of the path to find the closest one with a Makefile
	testPath := b.testPath
	for lcp != "." {
		testPath = path.Join(b.testPath, lcp)
		makefile := path.Join(testPath, "Makefile")
		fmt.Printf("using '%v' as base makefile\n", makefile)
		if s, err := os.Stat(makefile); err == nil && !s.IsDir() {
			break
		}
		lcp = filepath.Dir(lcp)

	}
	contents := []byte("success")
	out, origErr := exec.Command("make", "-C", testPath, "prof_test").CombinedOutput()
	if origErr != nil {
		contents = []byte("failure")
	}

	fmt.Printf("%v\n", string(out))
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
