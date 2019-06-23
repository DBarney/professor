# Professor
A simple utility that allows local runs of tests to be reported to branches as evidence that the tests do in fact pass.

This was built beacuse I get really tired of waiting for tests to get triggered, build, upload, and then set the status. Normally by the time I have pushed a branch, I have already ran the tests to see that they work.

Because this is designed to be run on the same server, it can take advantage of caching results of tests, builds, etc. it should be faster then starting from scratch each time.

usage:
```
prof # runs in CI/CD mode. watches local branches for changes and runs makefile targets
prof {branch|tag|commit} # runs a single build and uploads the results if needed.
```