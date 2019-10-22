package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/dbarney/professor/internal/publisher"
	"github.com/dbarney/professor/types"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var repoMatch = regexp.MustCompile("git@([^:]+):([^/]+)/([^.]+)[.]git")

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

	output, err := exec.Command("git", "config", "--get", "remote.origin.url").CombinedOutput()
	origin := string(output)
	host := "github.com"
	owner := "someone"
	name := "something"

	if strings.HasPrefix(origin, "git@") || strings.HasPrefix(origin, "https://") {
		parts := repoMatch.FindStringSubmatch(origin)
		host = parts[1]
		owner = parts[2]
		name = parts[3]
	} else {
		fmt.Printf("unsupported git remote type\n")
		os.Exit(1)
	}

	if host != "github.com" {
		url := &url.URL{Scheme: "https", Host: host, Path: "/api/v3/"}
		client.BaseURL = url
		client.UploadURL = url
	}

	err = createStatus(client, owner, name, types.Pending, ref, nil)
	if err != nil {
		fmt.Printf("unable to set status to pending %v\n", err)
	}
	base, err := os.Getwd()
	if err != nil {
		fmt.Printf("unable to determine working dir: %v\n", err)
		createStatus(client, owner, name, types.Error, ref, nil)
		os.Exit(1)
	}
	worktree := path.Join(base, ".git", "professor")

	os.RemoveAll(worktree)
	err = os.MkdirAll(worktree, 0777)
	if err != nil {
		fmt.Printf("unable to make path: %v\n", err)
		createStatus(client, owner, name, types.Error, ref, nil)
		os.Exit(1)
	}
	cmd := exec.Command("git", "worktree", "add", worktree, ref)

	o, err := cmd.CombinedOutput()
	fmt.Printf("create tree: %v %v\n", string(o), err)
	if err != nil {
		fmt.Printf("unable to create worktree %v\n", err)
		createStatus(client, owner, name, types.Error, ref, nil)
		os.Exit(1)
	}
	body := &bytes.Buffer{}
	writer := io.MultiWriter(os.Stdout, body)

	start := time.Now()
	cmd = exec.Command(command[0], command[1:]...)
	cmd.Dir = worktree
	cmd.Stdout = writer
	cmd.Stderr = writer
	err = cmd.Run()
	stop := time.Now()

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

	file, err := publisher.WrapWithMarkdown(publisher.Result{
		Output:   string(body.Bytes()),
		Duration: stop.Sub(start),
	}, result)

	if err != nil {
		fmt.Printf("unable to create markdown file\n")
		os.Exit(1)
	}

	files := map[github.GistFilename]github.GistFile{
		"BuildResults.md": github.GistFile{
			Content: github.String(file),
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

	err = createStatus(client, owner, name, result, ref, gist.HTMLURL)
	if err != nil {
		fmt.Printf("unable to set final status %v\n", err)
		os.Exit(1)
	}
}

func createStatus(client *github.Client, owner, name string, status types.Status, sha string, url *string) error {
	repoStatus := github.RepoStatus{
		State:       github.String(status.String()),
		Description: github.String(status.Description()),
		Context:     github.String("professor/local-integration"),
		TargetURL:   url,
	}
	ctx := context.Background()
	_, _, err := client.Repositories.CreateStatus(ctx, owner, name, sha, &repoStatus)
	return err
}
