# VCS Mock Proxy
_a.k.a. VCS Moxy_

It is written in Go and relies on the HTTP intercept capabilities of [ICAP](https://tools.ietf.org/html/rfc3507), as implemented in [go-icap/icap](https://github.com/go-icap/icap). In short, ICAP allows us to specify a set of criteria to match all requests against and then routes them accordingly. When a request hits the proxy, if it matches this criteria then a semi-hardcoded response is automatically short circuited in. If it does not match the criteria, the request proceeds as normal.

## Disclaimer

This is not an officially supported HashiCorp product.

## Layout

In general, this project attempts to stick to the layout prescribed in [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

`build`: CircleCI and Docker build scripts. Also includes relevant squid proxy [config](build/package/docker/configs/squid.conf) and [setup script](build/package/docker/scripts/squid-icap-init.sh).

`cmd`: Contains main.go file for building vcs-mock-proxy.

`deployments`: Configs for publishing builds of vcs-mock-proxy.

`hack`: Bash scripts for running vcs-mock-proxy locally, since a straight `docker-compose up` won't work.

`mocks`: Faux endpoints for testing the proxy redirect within this project. The `atlas` project houses endpoints used by itself in [atlas/integration-tests-api](https://github.com/hashicorp/atlas/tree/master/integration-tests-api/mocks).

`pkg`: The real meat of this whole setup.

## Getting started

To develop locally, run `./hack/local-dev-up.sh` to start the proxy and client containers. The result of this script will leave you in a bash shell inside the client container. From here, you can [create substitution variables](pkg/mock/mock.go#L199) and test API endpoints you've hardcoded via a GET curl.

To keep an eye on the proxy logs, including everyone's favorite debugger `Printf`, tail the proxy container logs via `docker logs -f deployments_proxy_1` (or whatever name your proxy's container happens to have).

## Dependency Management

This project uses Go [modules](https://github.com/golang/go/wiki/Modules) without a vendor dir.

This Wiki will probably have better advice about adding / upgrading dependencies than can be stated here.

## Routes File

The Routes file defines which endpoints should be mocked. It is defined in HCL, and should be stored in your mocks directory at `routes.hcl`.

Create each route as a block. The host should match the hostname of the site you want to mock. The path matches the path, but you can use rails route style `:foo` substitutions to make dynamic URLs. The type can either be `git` or `http`, depending on whether you want to mock a git clone operation or not.

```hcl
# You can reach this at any /orgs/$VALUE/repos request, and a substitution will
# be added that replaces {key=org, value=$VALUE}
route {
    host = "api.github.com"
    path = "/orgs/:org/repos"
    type = "http"
}
```

Do not create overlapping routes. This will cause an error, as the mock routing logic cannot determine which route to apply to a given request.

## Mocking Git Clones

vcs-mock-proxy also supports mocking Git Clones made via HTTP. To do so, add a route to your routes.hcl file:

```hcl
route {
    host = "github.com"
    path = "/example-repo"
    type = "git"
}
```

You'll next need to add a directory (for this example) at `/mocks/git/github.com/example-repo`. In this repo, you can then run a script `/hack/prep-git-mocks.sh` to automatically initialize a git repo. Don't commit this repo, as managing submodules is a pain. You can unstage that initialization with another script `/hack/unstage-git-mocks.sh`.

Once you've added a route, a directory exists at the correct path, and you've initialized git in it, you can run a clone by starting up the local dev environment and making an HTTP clone request:

```
./hack/local-dev-up.sh
git clone http://github.com/example-repo
```

## Mocking Different Response Codes

By default, all mocks return a 200 when they succeed. That's not the only possible thing you might need to mock though. In order to request a different response code along for a mock, send along the header `X-Desired-Response-Code`.

```
./hack/local-dev-up.sh
curl --head --header "X-Desired-Response-Code: 204" example.com
```

## SSL Certificates and Mocking HTTPS requests

Using Squid's SSL Bump configuration, VCS Mock Proxy can also act as an `https_proxy` and successfully mock upstream requests to HTTPS endpoints.

It does so in this local configuration using a self-signed certificate in `/certs`. This self-signed cert is automatically trusted for local dev, and shoud work out of the box.

If configuring VCS Mock Proxy in another environment, you will need to volume mount a self-signed certificate to `/etc/squid/ssl_cert/ca.pem`, and trust that certificate on any system attempting to use VCS Mock Proxy as an `https_proxy`.

To generate these certificates, use the script: `/hack/gen-certs.sh`. This may also be useful example code if you need to incorporate self-signed certs into another system using VCS Mock Proxy.
