package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var ErrNotFound = fmt.Errorf("the status was not found")
var ErrNoMakefile = fmt.Errorf("no makefile was found")

type Builder struct {
	clone     *git.Repository
	makefile  string
	buildPath string
	testPath  string
}

func NewBuilder(clone *git.Repository, makefile, buildPath, testPath string) *Builder {
	return &Builder{
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
	err = tree.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(sha),
		Force: true,
	})
	if err != nil {
		return err
	}
	if _, err := os.Stat(b.makefile); err != nil {
		err = ioutil.WriteFile(statusPath, []byte("error"), 0666)
		if err != nil {
			return err
		}
		return ErrNoMakefile
	}
	contents := []byte("success")
	out, origErr := exec.Command("make", "-C", b.testPath, "prof_test").CombinedOutput()
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
