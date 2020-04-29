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
