# Professor
A simple utility that allows local runs of tests to be reported to branches as evidence that the tests do in fact pass.

This was built beacuse I get really tired of waiting for tests to get triggered, build, upload, and then set the status. Normally by the time I have pushed a branch, I have already ran the tests to see that they work.

Because this is designed to be run on the same server, it can take advantage of caching results of tests, builds, etc. it should be faster then starting from scratch each time.

### Usage
```
prof # runs in CI/CD mode. watches local branches for changes and runs makefile targets
prof {branch|tag|commit} # runs a single build and uploads the results.
prof --comand 'make build' # override the command used to build
prof --origin 'git@github.com:dbarney/professor' # pull changes from this remote and don't try and use the current folder as a source for changes.
prof --auto-publish # don't wait for a remote to be updated, after the build works, upload the result
prof --build remote/origin # build refs matching this pattern
```

### Example builds
[Failure](https://gist.github.com/DBarney/d1e7920fcf6ae484d397430c1febea06)

[Success](https://gist.github.com/DBarney/61e0f6068911f125dc377600e642290a)

### configuration
Currently the only configuration that needs to happen is the Enviroment Variable `PROFESSOR_TOKEN` needs to be set to the value of a personal access token in github with the `gist` and `repo` permissions. This is used to update statuses, and to upload gists.

A lot of other settings currently aren't exposed and are set by reading the git config, and by setting fairly sane defaults. That being said, if need be in the future this can be changed with a simple PR to make them more configurable.

### future ideas?
```
prof --limit-to-users 'dbarney,john' # ony run commits by these users
prof --command 'make build-production' '--release './bin' # maybe create a github release? tag based maybe?
```