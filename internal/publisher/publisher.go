package publisher

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dbarney/professor/internal/storage"
	"github.com/dbarney/professor/types"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Publisher holds the configuration for publishing results
// to a gist on github
type Publisher struct {
	store  *storage.Store
	client *github.Client
	owner  string
	name   string
}

// NewPublisher creates and configures a publisher
func NewPublisher(host string, store *storage.Store, token, owner, name string) *Publisher {
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
		store:  store,
		client: client,
		owner:  owner,
		name:   name,
	}
}

func (p *Publisher) Publish(sha string) error {
	once := true
	for {
		collection, err := p.store.Get(sha)
		if err != nil {
			return err
		}
		elem := collection.Last(types.StatusOnly)
		if elem == nil {
			return fmt.Errorf("no final status was found")
		}
		event := elem.(*types.Event)
		if event.Status == types.Pending {
			if once {
				once = false
				err = p.createStatus(types.Pending, sha, nil)
				if err != nil {
					return err
				}
			}
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	return p.sendFinalStatus(sha)
}

func (p *Publisher) sendFinalStatus(sha string) error {
	collection, err := p.store.Get(sha)
	if err != nil {
		return err
	}

	lastElem := collection.Last(types.StatusOnly)
	firstElem := collection.First(types.StatusOnly)
	lastEvent := lastElem.(*types.Event)
	firstEvent := firstElem.(*types.Event)
	logs := collection.Take(types.LogOnly)

	body := []byte{}
	for _, elem := range logs {
		log := elem.(*types.Event)
		body = append(body, '\n')
		body = append(body, log.Data...)
	}
	body = body[1:]

	file, err := wrapWithMarkdown(Result{
		Output:   string(body),
		Duration: lastEvent.Time.Sub(firstEvent.Time),
	}, lastEvent.Status)

	if err != nil {
		return err
	}
	files := map[github.GistFilename]github.GistFile{
		"BuildResults.md": github.GistFile{
			Content: github.String(file),
		},
	}

	description := fmt.Sprintf("Professor Build for %v: %v", sha, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
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
	return p.createStatus(lastEvent.Status, sha, gist.HTMLURL)
}

func (p *Publisher) createStatus(status types.Status, sha string, url *string) error {
	repoStatus := github.RepoStatus{
		State:       github.String(status.String()),
		Description: github.String(status.Description()),
		Context:     github.String("professor/local-integration"),
		TargetURL:   url,
	}
	ctx := context.Background()
	_, _, err := p.client.Repositories.CreateStatus(ctx, p.owner, p.name, sha, &repoStatus)
	return err
}
