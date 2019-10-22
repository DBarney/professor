package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/dbarney/professor/types"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: prof {ref} {command}")
		os.Exit(1)
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Printf("GITHUB_TOKEN was empty\n")
		os.Exit(1)
	}
	user := os.Getenv("GITHUB_USER")
	if user == "" {
		fmt.Printf("GITHUB_USER was empty\n")
		os.Exit(1)
	}

	ref := os.Args[1]
	command := os.Args[2:]
	fmt.Printf("%v %v\n", ref, command)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// once we deal with remotes not stored on github, we will get
	// this working again
	/*if host != "github.com" {
		url := &url.URL{Scheme: "https", Host: host, Path: "/api/v3/"}
		client.BaseURL = url
		client.UploadURL = url
	}*/

	err := createStatus(client, types.Pending, ref, nil)
	if err != nil {
		fmt.Printf("unable to set status to pending %v\n", err)
	}
	base, err := os.Getwd()
	if err != nil {
		fmt.Printf("unable to determine working dir: %v\n", err)
		createStatus(client, types.Error, ref, nil)
		os.Exit(1)
	}
	worktree := path.Join(base, ".git", "professor")

	os.RemoveAll(worktree)
	err = os.MkdirAll(worktree, 0777)
	if err != nil {
		fmt.Printf("unable to make path: %v\n", err)
		createStatus(client, types.Error, ref, nil)
		os.Exit(1)
	}
	cmd := exec.Command("git", "worktree", "add", worktree, ref)

	o, err := cmd.CombinedOutput()
	fmt.Printf("create tree: %v %v\n", string(o), err)
	if err != nil {
		fmt.Printf("unable to create worktree %v\n", err)
		createStatus(client, types.Error, ref, nil)
		os.Exit(1)
	}

	cmd = exec.Command(command[0], command[1:]...)
	cmd.Dir = worktree
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	result := types.Success
	if err != nil {
		fmt.Printf("command failed %v\n", err)
		result = types.Failure
	}

	err = os.RemoveAll(worktree)
	if err != nil {
		fmt.Printf("clean up failed: %v\n", err)
	}
	cmd = exec.Command("git", "worktree", "prune")

	o, err = cmd.CombinedOutput()
	fmt.Printf("prune tree: %v %v\n", string(o), err)

	fmt.Printf("reporting status\n")

	files := map[github.GistFilename]github.GistFile{
		"BuildResults.md": github.GistFile{
			Content: github.String("just a test"),
		},
	}

	description := fmt.Sprintf("Professor Build for %v: %v", ref, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	gist := &github.Gist{
		Description: github.String(description),
		Public:      github.Bool(true),
		Files:       files,
	}
	gist, _, err = client.Gists.Create(ctx, gist)
	if err != nil {
		fmt.Printf("unable to create gist %v\n", err)
		os.Exit(1)
	}

	err = createStatus(client, result, ref, gist.HTMLURL)
	if err != nil {
		fmt.Printf("unable to set final status %v\n", err)
		os.Exit(1)
	}
}

func createStatus(client *github.Client, status types.Status, sha string, url *string) error {
	repoStatus := github.RepoStatus{
		State:       github.String(status.String()),
		Description: github.String(status.Description()),
		Context:     github.String("professor/local-integration"),
		TargetURL:   url,
	}
	ctx := context.Background()
	_, _, err := client.Repositories.CreateStatus(ctx, "dbarney", "professor", sha, &repoStatus)
	return err
}
