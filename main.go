package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dbarney/professor/internal/api"
	"github.com/dbarney/professor/internal/builder"
	"github.com/dbarney/professor/internal/publisher"
	"github.com/dbarney/professor/internal/repo"

	"github.com/logrusorgru/aurora"
	"gopkg.in/src-d/go-git.v4"
	git_config "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type flags struct {
	target      string
	origin      string
	autoPublish bool
	build       string
	check       time.Duration
	makefile    string
}

func main() {
	flags := &flags{}
	flag.StringVar(&flags.target, "target", "test", "the target that should be built")
	flag.StringVar(&flags.makefile, "makefile", "", "a path to a folder where a makefile should be used (the root of the repo is: '/')")
	flag.StringVar(&flags.origin, "origin", "", "the remote to use as the origin, defaults to the local directory")
	flag.BoolVar(&flags.autoPublish, "auto-publish", false, "trigger publishing when builds finish")
	flag.StringVar(&flags.build, "build", "heads/*", "the refs to monitor to trigger builds")
	flag.DurationVar(&flags.check, "check-every", time.Duration(0), "how often to poll to changes in the origin remote")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		// we run in headless mode, building everything that changes
		headlessRun(flags)
	} else if len(args) == 1 {
		// try to resolve the ref to a commit and only build and publish that.
		// should we support git style references? @~2 etc?
		singleRun(flags, args[0])
	} else {
		fmt.Printf("usage: prof {ref|sha|tag|branch}")
		os.Exit(1)
	}
}

func singleRun(flags *flags, arg string) {
	fmt.Printf("running single build: %v\n", arg)
	config, err := getConfig(flags.origin)
	if err != nil {
		panic(err)
	}

	fmt.Printf("opening repo.\n")
	original, err := git.PlainOpen(config.topLevel)
	if err != nil {
		panic(err)
	}

	sha, err := argToSha(original, arg)
	if err != nil {
		panic(err)
	}

	build := builder.NewBuilder(original, flags.makefile, flags.target, config.buildPath, config.workingPath)

	pub := publisher.NewPublisher(config.host, build, config.token, config.owner, config.name)

	err = build.Build(sha)
	switch err {
	case nil:
		fmt.Println(aurora.Green("build was sucessful!"))
	case builder.ErrNoMakefile:
		fmt.Println("no Makefile was found skipping tests.")
		os.Exit(1)
	case builder.ErrNoChanges:
		fmt.Println(aurora.Yellow("no changes were detected."))
		os.Exit(0)
	default:
		fmt.Printf("%v,", err)
		fmt.Printf("%v %v\n", aurora.Red("build Failed"), err)
	}
	res, err := build.GetResults(sha)
	if err == nil {
		fmt.Printf("\noutput: \n%v\n", res["BuildResults.txt"])
	}

	err = pub.Publish(sha)
	switch err {
	case nil:
		fmt.Println(aurora.Green("sucessfully published build results."))
	default:
		fmt.Printf("%v %v", aurora.Red("unable to publish results:"), err)
		os.Exit(1)
	}

}

func headlessRun(flags *flags) {
	config, err := getConfig(flags.origin)
	if err != nil {
		panic(err)
	}
	// watch for changes on local branches and remote branches
	repository := repo.New(config.gitFolder)
	local, err := repository.WatchLocalBranches(flags.build)
	if err != nil {
		panic(err)
	}
	var source chan *repo.BranchEvent
	var remote <-chan *repo.BranchEvent
	if !flags.autoPublish {
		remote, err = repository.WatchRemoteBranches()
		if err != nil {
			panic(err)
		}
	} else {
		source = make(chan *repo.BranchEvent)
		remote = source
	}

	fmt.Printf("opening repo.\n")
	original, err := git.PlainOpen(config.topLevel)
	if err != nil {
		panic(err)
	}
	if flags.check != 0 {
		ticker := time.NewTicker(flags.check)
		defer ticker.Stop()
		refspec := buildRefSepc(flags.build)
		go func() {
			for range ticker.C {
				fmt.Println("checking for changes")
				err := original.Fetch(&git.FetchOptions{
					RefSpecs: []git_config.RefSpec{git_config.RefSpec(refspec)},
					Auth: &http.BasicAuth{
						Username: "dbarney",
						Password: config.token,
					}})
				fmt.Printf("fetched changes %v\n", err)
			}
		}()
	}

	// start the API
	s := make(chan *api.Set)
	api.Run(s)

	// start the build process
	build := builder.NewBuilder(original, flags.makefile, flags.target, config.buildPath, config.workingPath)
	stream := build.GetStream()
	go handleLocalChanges(local, build, source)

	go func() {
		var events chan *api.Event
		for stream := range stream {
			switch stream.T {
			case builder.Start:
				events = make(chan *api.Event)
				set := &api.Set{
					Name:   stream.Message,
					Events: events,
				}
				s <- set
			case builder.Stop:
				close(events)
			case builder.Message:
				events <- &api.Event{
					Name: "output",
					Data: stream.Message,
				}
			case builder.Target:
				events <- &api.Event{
					Name: "target",
					Data: stream.Message,
				}
			}
		}
	}()
	// start the reporting process
	pub := publisher.NewPublisher(config.host, build, config.token, config.owner, config.name)
	handleRemoteChanges(remote, pub)
}

func handleRemoteChanges(changes <-chan *repo.BranchEvent, pub *publisher.Publisher) {
	for c := range changes {
		fmt.Printf("detected a remote being updated %v.\n", c.SHA)
		go func(sha string) {
			err := pub.Publish(sha)
			if os.IsNotExist(err) {
				return
			} else if err == nil {
				fmt.Println(aurora.Green("sucessfully published build results."))
			} else {
				fmt.Printf("%v %v\n", aurora.Red("unable to publish results:"), err)
			}
		}(c.SHA)
	}
}

func handleLocalChanges(changes <-chan *repo.BranchEvent, build *builder.Builder, next chan<- *repo.BranchEvent) {
	for c := range changes {
		fmt.Printf("detected a local branch being updated, building %v\n", c.SHA)
		err := build.Build(c.SHA)
		switch err {
		case nil:
			fmt.Println(aurora.Green("build was sucessful!"))
		case builder.ErrNoMakefile:
			fmt.Println("no Makefile was found skipping tests.")
			continue
		case builder.ErrNoChanges:
			fmt.Println(aurora.Yellow("no changes were detected."))
			continue
		default:
			fmt.Printf("%v %v\n", aurora.Red("build Failed."), err)
		}
		// forward the trigger if there is a next step
		if next != nil {
			next <- c
		}
	}
}

func buildRefSepc(build string) string {
	switch {
	case strings.HasPrefix(build, "heads/"):
		build = strings.TrimPrefix(build, "heads/")
		return fmt.Sprintf("+refs/heads/%v:refs/remotes/origin/%v", build, build)
	case strings.HasPrefix(build, "remotes/origin/"):
		build = strings.TrimPrefix(build, "remotes/origin/")
		return fmt.Sprintf("+refs/heads/%v:refs/remotes/origin/%v", build, build)
	}
	// don't really know what to do with this
	return "+refs/heads/*:refs/remotes/origin/*"
}

// try and take a string and discover if it is something representing
// a sha in the git repositoy
func argToSha(repo *git.Repository, arg string) (string, error) {
	tag, err := repo.Tag(arg)
	if err == nil {
		return tag.Hash().String(), nil
	}
	// is it a local branch
	name := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%v", arg))
	ref, err := repo.Reference(name, true)
	if err == nil {
		return ref.Hash().String(), nil
	}
	// is it a remote branch
	name = plumbing.ReferenceName(fmt.Sprintf("refs/remotes/%v", arg))
	ref, err = repo.Reference(name, true)
	if err == nil {
		return ref.Hash().String(), nil
	}
	_, err = repo.CommitObject(plumbing.NewHash(arg))
	if err == nil {
		return arg, nil
	}
	return "", fmt.Errorf("%v was not a tag, a ref or a sha.\n", arg)
}
