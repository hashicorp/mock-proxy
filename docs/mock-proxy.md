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
`/certs`. This self-signed cert is automatically trusted for local dev.

To generate these certificates, use the script: `/hack/gen-certs.sh`. This may
also be useful example code if you need to incorporate self-signed certs into
another system using mock-proxy.

If configuring mock-proxy in another environment, you will need to volume mount
a self-signed certificate to `/etc/squid/ssl_cert/ca.pem`, and trust that
certificate on any system attempting to use mock-proxy as an `https_proxy`.

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
