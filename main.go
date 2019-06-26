package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dbarney/professor/internal/builder"
	"github.com/dbarney/professor/internal/publisher"
	"github.com/dbarney/professor/internal/repo"

	"github.com/logrusorgru/aurora"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type flags struct {
	command string
}

func main() {
	flags := &flags{}
	flag.StringVar(&flags.command, "command", "make test", "the command that should be run")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		// we run in headless mode, building everything that changes
		headlessRun(flags)
	} else if len(args) == 1 {
		// try to resolve the ref to a commit and only build and publish that.
		// should we support git style references? @~2 etc?
		singleRun(flags, args[1])
	} else {
		fmt.Printf("usage: prof {ref|sha|tag|branch}")
		os.Exit(1)
	}
}

func singleRun(flags *flags, arg string) {
	fmt.Printf("running single build: %v\n", arg)
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	repo, err := setupRepo(config)
	if err != nil {
		panic(err)
	}

	sha, err := argToSha(repo, arg)
	if err != nil {
		panic(err)
	}

	original, err := git.PlainOpen(config.topLevel)
	if err != nil {
		panic(err)
	}
	build := builder.NewBuilder(original, repo, flags.command, config.makefile, config.buildPath, config.testPath)

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
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	// watch for changes on local branches and remote branches
	repository := repo.New(config.topLevel)
	local, err := repository.WatchLocalBranches()
	if err != nil {
		panic(err)
	}
	remote, err := repository.WatchRemoteBranches()
	if err != nil {
		panic(err)
	}

	repo, err := setupRepo(config)
	if err != nil {
		panic(err)
	}

	original, err := git.PlainOpen(config.topLevel)
	if err != nil {
		panic(err)
	}
	// start the build process
	build := builder.NewBuilder(original, repo, flags.command, config.makefile, config.buildPath, config.testPath)
	go handleLocalChanges(local, build)

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

func handleLocalChanges(changes <-chan *repo.BranchEvent, build *builder.Builder) {
	for c := range changes {
		fmt.Printf("detected a local branch being updated, building %v\n", c.SHA)
		switch build.Build(c.SHA) {
		case nil:
			fmt.Println(aurora.Green("build was sucessful!"))
		case builder.ErrNoMakefile:
			fmt.Println("no Makefile was found skipping tests.")
		case builder.ErrNoChanges:
			fmt.Println(aurora.Yellow("no changes were detected."))
		default:
			fmt.Println(aurora.Red("build Failed."))
		}
	}
}

func setupRepo(config *config) (*git.Repository, error) {
	_, err := os.Stat(config.testPath)
	if os.IsNotExist(err) {
		fmt.Printf("cloning repo.\n")
		err = os.MkdirAll(config.workingPath, 0777)
		if err != nil {
			return nil, err
		}
		return git.PlainClone(config.testPath, false, &git.CloneOptions{
			URL: config.topLevel,
		})
	} else if err != nil {
		return nil, err
	} else {
		fmt.Printf("opening repo.\n")
		return git.PlainOpen(config.testPath)
	}
}

// try and take a string a descover if it is something representing
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
