package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

type config struct {
	topLevel    string
	testPath    string
	makefile    string
	buildPath   string
	workingPath string

	token string
	host  string
	owner string
	name  string
}

var urlMatch = regexp.MustCompile("git@([^:]+):([^/]+)/([^.]+)[.]git")
var pathMatch = regexp.MustCompile("([^/]+)/(.+)")

func getConfig() (*config, error) {
	// find the top level folder
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, err
	}

	topLevel := strings.TrimSpace(string(out))

	token := os.Getenv("PROFESSOR_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("PROFESSOR_TOKEN was empty")
	}

	base, err := git.PlainOpen(topLevel)
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

	return &config{
		topLevel:    topLevel,
		testPath:    path.Join(topLevel, ".git", "professor", "working"),
		makefile:    path.Join(topLevel, ".git", "professor", "working", "Makefile"),
		buildPath:   path.Join(topLevel, ".git", "professor", "builds"),
		workingPath: path.Join(topLevel, ".git", "professor"),

		token: token,
		host:  parts[1],
		owner: parts[2],
		name:  parts[3],
	}, nil
}
