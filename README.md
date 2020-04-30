# mock-proxy
_a.k.a. Moxie_

mock-proxy (a.k.a Moxie—short for mo(ck)(pro)xy) is a replacement proxy relying
on the HTTP intercept capabilities of [ICAP](https://tools.ietf.org/html/rfc3507),
as implemented in [go-icap/icap](https://github.com/go-icap/icap). This
replacement proxy is intended for writing integration tests that mock responses
from external services at the network level. In order to use mock-proxy, a test
environment specifies it as an HTTP proxy using `http_proxy` type environment
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

## Features

* Selectively mock endpoints, allowing some requests to hit the internet, and
others to be faked locally.
* Configure mocked routes using an HCL2 based Routes file.
* Dynamic URL support, allowing mocking traditional RESTful APIs easily.
* Templated responses using "Transformer" interface, including the built in
substitution-variables endpoint.
* Mock `git clone` operations using the `git` route type.
* Mock HTTPS endpoints using Squid's SSLBump feature.

See documentation and examples for more information.

## Disclaimer

This is not an officially supported HashiCorp product.

## Documentation

There is documentation of how to use mock-proxy features such as the Routes
file in the [docs](/docs) directory.

## Examples

See [examples of different use cases for mock-proxy](/docs/examples)

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

The local development and demonstration environment for mock-proxy relies on
Docker to orchestrate the various services (client, ICAP server, Squid server)
required to successfully mock requests. If you do not have a functional Docker
environment, see the
[getting started docs on Docker's website](https://docs.docker.com/get-started/#set-up-your-docker-environment).

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
