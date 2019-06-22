# Professor
A simple utility that allows local runs of tests to be reported to branches as evidence that the tests do in fact pass.


usage:
```

// against the head of MASTER
prof httpd://github.com/DBarney/lever
// against the head of a PR
prof https://github.com/DBarney/lever/pull/2
// against a specific commit
prof https://github.com/DBarney/lever/commit/534888e6792082f94ac59b00ceb867b88d78237f
// against a specific branch
prof https://github.com/DBarney/lever/tree/feature/stream-rewrite

the recomended way is to specify the defaults in a `.prof` file at the root of the repo:
url https://github.com/DBarney/lever

the the commands become shortened to:
prof
prof pull/2
prof commit/534888e6792082f94ac59b00ceb867b88d78237f
prof tree/feature/stream-rewrite

// watch for changes, and when new commits are added to the branch trigger a new test run
prof watch https://github.com/DBarney/lever/pull/2

// watch the repo for changes by my user and run tests against those
prof httpd://github.com/DBarney/lever


reporting status to central server
```