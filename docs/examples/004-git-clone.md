## Git Clone

London is working on a service that needs to mock out cloning git repositories.

To that end she generates a Routes file:

```hcl
route {
  host = "github.com"
  path = "/example-repo"
  type = "git"
}
```

And initializes a git repository at `/mocks/git/github.com/example-repo`.

Now she can clone the fake git repository with http clones:

```bash
export http_proxy=http://squid.proxy:8888
export https_proxy=http://squid.proxy:8888

git clone https://github.com/example-repo --depth=1
### Cloning into 'example-repo'...
### Unpacking objects: 100% (3/3), done.
```
