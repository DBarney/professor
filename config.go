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

	token string
	host  string
	owner string
	name  string
}

var urlMatch = regexp.MustCompile("git@([^:]+):([^/]+)/([^.]+)[.]git")
var pathMatch = regexp.MustCompile("([^/]+)/(.+)")

func getConfig(origin string) (*config, error) {
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
