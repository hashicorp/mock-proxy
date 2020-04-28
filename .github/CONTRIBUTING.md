# Contributing to mock-proxy

mock-proxy is an open source project and we appreciate contributions of various
kinds, including bug reports and fixes, enhancement proposals, documentation
updates, and user experience feedback. However, this is not an officially
supported HashiCorp product. We _are_ grateful for all contributions, but this
repository is primarily maintained by a small team at HashiCorp outside of
their main job responsibilities and is developed on a volunteer basis.

To record a bug report, enhancement proposal, or give any other product
feedback, please [open a GitHub issue](https://github.com/hashicorp/mock-proxy/issues/new).

**All communication on GitHub, the community forum, and other HashiCorp-provided
communication channels is subject to
[the HashiCorp community guidelines](https://www.hashicorp.com/community-guidelines).**

## Local Development Environment

mock-proxy has a simple proxying setup configured in docker-compose for use
while developing changes locally. To start up this local development
environment run:

```
./hack/local-dev-up.sh
```

This starts the Squid proxy, an ICAP server implemented by mock-proxy, and a
client that is using that proxy. You can see a simple mock implemented by
running:

```
curl https://example.com
```

## Running Tests

mock-proxy has a test suite implemented for its Golang components. To run the
tests, run:

```
go test ./...
```

For additional details run:

```
go test -v -race ./...
```

## External Dependencies

mock-proxy uses Go Modules for dependency management.

Our dependency licensing policy for mock-proxy excludes proprietary licenses
and "copyleft"-style licenses. We accept the common Mozilla Public License v2,
MIT License, and BSD licenses. We will consider other open source licenses
in similar spirit to those three, but if you plan to include such a dependency
in a contribution we'd recommend opening a GitHub issue first to discuss what
you intend to implement and what dependencies it will require so that the
mock-proxy team can review the relevant licenses for whether they meet our
licensing needs.

If you need to add a new dependency to mock-proxy or update the selected version
for an existing one, use `go get` from the root of the mock-proxy repository
as follows:

```
go get github.com/hashicorp/hcl/v2@2.0.0
```

This command will download the requested version (2.0.0 in the above example)
and record that version selection in the `go.mod` file. It will also record
checksums for the module in the `go.sum`.

To complete the dependency change and clean up any redundancy in the module
metadata files by running the following commands:

```
go mod tidy
```

Because dependency changes affect a shared, top-level file, they are more likely
than some other change types to become conflicted with other proposed changes
during the code review process. For that reason, and to make dependency changes
more visible in the change history, we prefer to record dependency changes as
separate commits that include only the results of the above commands and the
minimal set of changes to mock-proxy's own code for compatibility with the
new version:

```
git add go.mod go.sum
git commit -m "modules: go get github.com/hashicorp/hcl/v2@2.0.0"
```

You can then make use of the new or updated dependency in new code added in
subsequent commits.
