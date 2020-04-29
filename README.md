# mock-proxy
_a.k.a. Moxie_

mock-proxy (a.k.a Moxie—short for mo(ck)(pro)xy) is a replacement proxy relying
on the HTTP intercept capabilities of [ICAP](https://tools.ietf.org/html/rfc3507),
as implemented in [go-icap/icap](https://github.com/go-icap/icap). This
replacement proxy is intended for writing integration tests that mock responses
from external services at the network level. In order to use mock-proxy, a test
environment specificies it as an HTTP proxy using `http_proxy` type environment
variables (or however proxies are configured in your scenario). Because it
works at the network level and can be used by many services at once using proxy
config, it works well in microservices environments where many services need to
mock the same 3rd party dependencies.

In short, ICAP allows us to specify a set of criteria to match all requests
against and then route them accordingly. When a request hits the proxy, if it
matches this criteria then a semi-hardcoded response is automatically short
circuited in. If it does not match the criteria, the request proceeds as
normal.

![moxie flow diagram](/docs/images/mock-proxy-diagram.png)

## Disclaimer

This is not an officially supported HashiCorp product.

## Layout

In general, this project attempts to stick to the layout prescribed in [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

`build`: CircleCI and Docker build scripts. Also includes relevant squid proxy
[config](build/package/docker/configs/squid.conf) and
[setup script](build/package/docker/scripts/squid-icap-init.sh).

`certs`: Self signed certificates used to configure SSL Bump for local
development.

`cmd`: Contains main.go file for building mock-proxy.

`deployments`: Configs for publishing builds of mock-proxy.

`hack`: Bash scripts for assorted tasks such as running mock-proxy locally

`mocks`: Faux endpoints for testing the proxy redirect within this project. 

`pkg`: The main Go code directory which houses the implementation of the custom
ICAP server.

## Getting started

To develop locally, run `./hack/local-dev-up.sh` to start the proxy and client
containers. The result of this script will leave you in a bash shell inside the
client container. From here, you can create substitution variables and test API
endpoints you've hardcoded via a GET curl.

To keep an eye on the proxy logs, tail the proxy container logs via
`docker logs -f deployments_proxy_1` (or whatever name your proxy's container
happens to have).

## Dependency Management

This project uses Go [modules](https://github.com/golang/go/wiki/Modules)
without a vendor dir.

This Modules Wiki will probably have better advice about adding / upgrading
dependencies than can be stated here.

## Routes File

The Routes file defines which endpoints should be mocked. It is defined in HCL,
and should be stored in your mocks directory at `routes.hcl`.

Create each route as a block. The host should match the hostname of the site
you want to mock. The path matches the path, but you can use rails route style
`:foo` substitutions to make dynamic URLs. The type can either be `git` or
`http`, depending on whether you want to mock a git clone operation or not.

```hcl
# You can reach this at any /orgs/$VALUE/repos request, and a substitution will
# be added that replaces {key=org, value=$VALUE}
route {
    host = "api.github.com"
    path = "/orgs/:org/repos"
    type = "http"
}
```

Do not create overlapping routes. This will cause an error, as the mock routing
logic cannot determine which route to apply to a given request.

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

## Mocking Different Response Codes

By default, all mocks return a 200 when they succeed. That's not the only
possible thing you might need to mock though. In order to request a different
response code along for a mock, send along the header
`X-Desired-Response-Code`.

```
./hack/local-dev-up.sh
curl --head --header "X-Desired-Response-Code: 204" example.com
```

## SSL Certificates and Mocking HTTPS requests

Using Squid's SSL Bump configuration, mock-proxy can also act as an
`https_proxy` and successfully mock upstream requests to HTTPS endpoints.

It does so in this local configuration using a self-signed certificate in
`/certs`. This self-signed cert is automatically trusted for local dev, and
should work out of the box.

If configuring mock-proxy in another environment, you will need to volume mount
a self-signed certificate to `/etc/squid/ssl_cert/ca.pem`, and trust that
certificate on any system attempting to use mock-proxy as an `https_proxy`.

To generate these certificates, use the script: `/hack/gen-certs.sh`. This may
also be useful example code if you need to incorporate self-signed certs into
another system using mock-proxy.
