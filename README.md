# VCS Mock Proxy
_a.k.a. VCS Moxy_

`vcs-mock-proxy-poc` is a POC intercept proxy for VCS tests in [atlas](https://github.com/hashicorp/atlas). There is an RFC in progress to request it be allowed to drop the `-poc`.

It is written in Go and relies on the HTTP intercept capabilities of [ICAP](https://tools.ietf.org/html/rfc3507), as implemented in [go-icap/icap](https://github.com/go-icap/icap).

## Layout

In general, this project attempts to stick to the layout prescribed in [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

`build`: CircleCI and Docker build scripts. Also includes relevant squid proxy [config](build/package/docker/configs/squid.conf) and [setup script](build/package/docker/scripts/squid-icap-init.sh).

`cmd`: Contains main.go file for building vcs-mock-proxy.

`deployments`: Configs for publishing builds of vcs-mock-proxy.

`hack`: Bash scripts for running vcs-mock-proxy locally, since a straight `docker-compose up` won't work.

`mocks`: Faux endpoints for testing the proxy redirect within this project. The `atlas` project houses endpoints used by itself in [atlas/integration-tests-api](https://github.com/hashicorp/atlas/tree/master/integration-tests-api/mocks).

`pkg`: The real meat of this whole setup.

## Dependency Management

This project uses Go [modules](https://github.com/golang/go/wiki/Modules) without a vendor dir.

This Wiki will probably have better advice about adding / upgrading dependencies than can be stated here.
