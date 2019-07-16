package repo

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
	return l.watch(path.Join(l.path, "refs", "remotes", "origin"))
}

// Watch branch returns a channel that reports if something changes
// with a local branch
func (l *Local) WatchLocalBranches(refs string) (<-chan *BranchEvent, error) {
	refs = strings.TrimSuffix(refs, "/*")
	fmt.Printf("watching %v\n", refs)
	return l.watch(path.Join(l.path, "refs", refs))
}

func (l *Local) watch(path string) (<-chan *BranchEvent, error) {
	path = strings.TrimSuffix(path, "/*")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	l.watchers = append(l.watchers, watcher)

	foundPath := resolvePath(path)
	err = watcher.Add(foundPath)
	if err != nil {
		return nil, err
	}
	err = watchPath(watcher, foundPath, path, nil)
	if err != nil {
		return nil, err
	}
	results := make(chan *BranchEvent)
	go func() {
		prev := ""
		t := time.Time{}
		for {
			select {
			case event := <-watcher.Events:
				// everything above Write is not something we are
				// interested in.
				if event.Op > fsnotify.Write || strings.HasSuffix(event.Name, ".lock") {
					continue
				}

				stat, err := os.Stat(event.Name)
				if err != nil {
					// the item no longer exists, nothing to see here
					continue
				}
				if stat.IsDir() &&
					(strings.HasPrefix(event.Name, path) ||
						strings.HasPrefix(path, event.Name)) {
					fmt.Printf("got a new folder, %v %v\n", event.Name, path)
					watcher.Add(event.Name)
					err = watchPath(watcher, event.Name, path, watcher.Events)
					fmt.Printf("%v\n", err)
					continue
				}
				if !strings.HasPrefix(event.Name, path) {
					fmt.Printf("change wasn't in the right path\n")
					continue
				}

				sha, err := getSha(event.Name)
				if err != nil {
					panic(err)
				}
				if len(sha) != 40 || sha == prev && time.Now().Sub(t) < time.Second {
					continue
				}
				prev = sha
				t = time.Now()
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

func watchPath(watcher *fsnotify.Watcher, path string, prefix string, events chan fsnotify.Event) error {
	return filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasPrefix(prefix, dir) && !strings.HasPrefix(dir, prefix) {
			return nil
		}
		if info.IsDir() {
			return watcher.Add(dir)
		}
		if events == nil {
			return nil
		}
		events <- fsnotify.Event{
			Name: dir,
			Op:   fsnotify.Write,
		}
		return nil
	})
}

func resolvePath(path string) string {
	for {
		_, err := os.Stat(path)
		if err == nil {
			// this will always eventually work.
			// as the path will either be '.' or '/'
			break
		}
		path = filepath.Dir(path)
	}
	return path
}
func getSha(file string) (string, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}
