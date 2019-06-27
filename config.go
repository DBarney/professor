package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type config struct {
	topLevel    string
	testPath    string
	buildPath   string
	workingPath string
	gitFolder   string
	refspec     string

	token string
	host  string
	owner string
	name  string
}

var urlMatch = regexp.MustCompile("git@([^:]+):([^/]+)/([^.]+)[.]git")
var pathMatch = regexp.MustCompile("([^/]+)/(.+)")

func getConfig(origin, build string) (*config, error) {
	c := &config{}

	c.token = os.Getenv("PROFESSOR_TOKEN")
	if c.token == "" {
		return nil, fmt.Errorf("PROFESSOR_TOKEN was empty")
	}
	if strings.HasPrefix(origin, "git@") || strings.HasPrefix(origin, "https://") {
		parts := urlMatch.FindStringSubmatch(origin)
		if len(parts) != 4 {
			panic(fmt.Sprintf("remote didn't match %v", parts))
		}

		c.host = parts[1]
		c.owner = parts[2]
		c.name = parts[3]
		c.topLevel = "./"
		c.gitFolder = c.topLevel
		switch {
		case strings.HasPrefix(build, "remotes/origin/"):
			rest := strings.TrimPrefix(build, "remotes/origin/")
			// pretty standard, but this allows namespacing
			c.refspec = fmt.Sprintf("refs/heads/%v:refs/remote/origin/%v", rest, rest)
		case strings.HasPrefix(build, "heads/"):
			// maybe someone will push some local branches here, so we don't
			// do any fetching.
			c.refspec = ""
		case strings.HasPrefix(build, "tags/"):
			// hey we can just build tags as well!
			rest := strings.TrimPrefix(build, "tags/")
			c.refspec = fmt.Sprintf("refs/tags/%v:refs/tags/%v", rest, rest)
		}

		// this really needs to be a better check.
		if _, err := os.Stat("./HEAD"); err != nil {
			// uses https cloning because we have a token already.
			// don't need an ssh key this way
			url := fmt.Sprintf("https://%v/%v/%v.git", c.host, c.owner, c.name)
			// a remote that needs to be cloned down
			_, err := git.PlainClone(c.topLevel, true, &git.CloneOptions{
				URL: url,
				Auth: &http.BasicAuth{
					Username: "dbarney",
					Password: c.token,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		c.testPath = path.Join(c.topLevel)
		c.buildPath = path.Join(c.topLevel, "professor", "builds")
		c.workingPath = path.Join(c.topLevel, "professor")

	} else {
		// try to discover from the current working directory
		// find the top level folder
		command := exec.Command("git", "rev-parse", "--show-toplevel")
		if origin != "" {
			// if origin isn't nil, then lets use it as the path to check
			// for a git repo
			command.Dir = origin
		}
		out, err := command.Output()
		if err != nil {
			return nil, err
		}
		c.topLevel = strings.TrimSpace(string(out))

		c.gitFolder = path.Join(c.topLevel, ".git")
		c.workingPath = path.Join(c.gitFolder, "professor")
		c.testPath = path.Join(c.workingPath, "working")
		c.buildPath = path.Join(c.workingPath, "builds")
		switch {
		case strings.HasPrefix(build, "remotes/origin/"):
			rest := strings.TrimPrefix(build, "remotes/origin/")
			// our remote is a local disk copy which has its own remote/origin.
			// so we fetch those ones.
			c.refspec = fmt.Sprintf("refs/remote/origin/%v:refs/remote/origin/%v", rest, rest)
		case strings.HasPrefix(build, "heads/"):
			// building just the local refs
			rest := strings.TrimPrefix(build, "heads/")
			c.refspec = fmt.Sprintf("refs/heads/%v:refs/heads/%v", rest, rest)
		case strings.HasPrefix(build, "tags/"):
			// hey we can just build tags as well!
			rest := strings.TrimPrefix(build, "tags/")
			c.refspec = fmt.Sprintf("refs/tags/%v:refs/tags/%v", rest, rest)
		}

		base, err := git.PlainOpen(c.topLevel)
		if err != nil {
			return nil, err
		}

		origin, err := base.Remote("origin")
		if err != nil {
			return nil, err
		}

		remote := origin.Config().URLs[0]
		parts := urlMatch.FindStringSubmatch(remote)
		if len(parts) != 4 {
			panic("remote didn't match")
		}

		c.host = parts[1]
		c.owner = parts[2]
		c.name = parts[3]
	}
	return c, nil
}
