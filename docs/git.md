## Mocking Git Clones

mock-proxy also supports mocking Git Clones made via HTTP. To do so, add a
route to your routes.hcl file:

```hcl
route {
    host = "github.com"
    path = "/example-repo"
    type = "git"
}
```

You'll next need to add a directory (for this example) at
`/mocks/git/github.com/example-repo`. In this repo, you can then run a script
`/hack/prep-git-mocks.sh` to automatically initialize a git repo. Don't commit
this repo, as we don't want to manage git submodules here. You can unstage that
initialization with another script `/hack/unstage-git-mocks.sh`.

Once you've added a route, a directory exists at the correct path, and you've
initialized git in it, you can run a clone by starting up the local dev
environment and making an HTTP clone request:

```
./hack/local-dev-up.sh
git clone http://github.com/example-repo
```
