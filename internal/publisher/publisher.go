package publisher

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Publisher struct {
	storage Storage
	client  *github.Client
	owner   string
	name    string
}

// Storage represents how build results are stored and retreived
type Storage interface {
	GetStatus(sha string) (string, error)
	GetResults(sha string) (map[string]string, error)
}

func NewPublisher(host string, store Storage, token, owner, name string) *Publisher {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	if host != "github.com" {
		url := &url.URL{Scheme: "https", Host: host, Path: "/api/v3/"}
		client.BaseURL = url
		client.UploadURL = url
	}
	return &Publisher{
		storage: store,
		client:  client,
		owner:   owner,
		name:    name,
	}
}

func (p *Publisher) Publish(sha string) error {
	once := true
	for {
		status, err := p.storage.GetStatus(sha)
		if err != nil {
			return err
		}
		switch {
		case status == "pending" && once:
			once = false
			err = p.createStatus(sha, status, nil)
			if err != nil {
				return err
			}
			fallthrough
		case status == "pending":
			time.Sleep(time.Second * 5)
			continue
		}
		return p.sendFinalStatus(sha)
	}
	return nil
}

func (p *Publisher) sendFinalStatus(sha string) error {
	status, err := p.storage.GetStatus(sha)
	if err != nil {
		return err
	}
	results, err := p.storage.GetResults(sha)
	if err != nil {
		return err
	}

	files := map[github.GistFilename]github.GistFile{}
	for name, body := range results {
		files[github.GistFilename(name)] = github.GistFile{
			Content: github.String(body),
		}
	}
	description := fmt.Sprintf("Professor Build #(%v): %v", sha, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	gist := &github.Gist{
		Description: github.String(description),
		Public:      github.Bool(true),
		Files:       files,
	}
	ctx := context.Background()
	gist, _, err = p.client.Gists.Create(ctx, gist)
	if err != nil {
		return err
	}
	return p.createStatus(status, sha, gist.HTMLURL)
}

func (p *Publisher) createStatus(status, sha string, url *string) error {
	var message string
	switch status {
	case "pending":
		message = "the build is pending..."
	case "success":
		message = "the build was sucessful!"
	case "failure":
		message = "something went wrong."
	case "error":
		message = "the build failed."
	default:
		return fmt.Errorf("unknown build status %v", status)
	}
	repoStatus := github.RepoStatus{
		State:       github.String(status),
		Description: github.String(message),
		Context:     github.String("professor/local-integration"),
		TargetURL:   url,
	}
	ctx := context.Background()
	_, _, err := p.client.Repositories.CreateStatus(ctx, p.owner, p.name, sha, &repoStatus)
	return err
}
