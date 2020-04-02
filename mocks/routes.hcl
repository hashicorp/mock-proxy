route {
  host = "example.com"
  path = "/"
  type = "http"
}

# You can reach this at any /orgs/$VALUE/repos request, and a substitution will
# be added that replaces {key=org, value=$VALUE}
route {
  host = "api.github.com"
  path = "/orgs/:org/repos"
  type = "http"
}

route {
  host = "github.com"
  path = "/example-repo"
  type = "git"
}
