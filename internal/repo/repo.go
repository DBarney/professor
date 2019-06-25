package repo

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Local represents a local git repository.
type Local struct {
	path     string
	tracking map[string]bool
	watchers []*fsnotify.Watcher
}

// BranchEvent represents something happening in a repository
type BranchEvent struct {
	SHA string
}

// New creates and returns a new Local git repositoy
func New(path string) *Local {
	return &Local{
		path:     path,
		tracking: map[string]bool{},
	}
}

func (l *Local) WatchRemoteBranches() (<-chan *BranchEvent, error) {
	return l.watch(path.Join(l.path, ".git", "refs", "remotes", "origin"))
}

// Watch branch returns a channel that reports if something changes
// with a local branch
func (l *Local) WatchLocalBranches() (<-chan *BranchEvent, error) {
	return l.watch(path.Join(l.path, ".git", "refs", "heads"))
}

func (l *Local) watch(path string) (<-chan *BranchEvent, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	l.watchers = append(l.watchers, watcher)

	err = watcher.Add(path)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		return watcher.Add(path)
	})
	if err != nil {
		return nil, err
	}
	results := make(chan *BranchEvent)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Create || strings.HasSuffix(event.Name, ".lock") {
					continue
				}

				stat, err := os.Stat(event.Name)
				if err != nil {
					panic(err)
				}
				if stat.IsDir() {
					continue
				}
				body, err := ioutil.ReadFile(event.Name)
				if err != nil {
					panic(err)
				}
				sha := strings.TrimSpace(string(body))
				results <- &BranchEvent{
					SHA: sha,
				}
			case err := <-watcher.Errors:
				panic(err)
			}
		}
	}()
	return results, nil
}
