route {
    host = "example.com"
    path = "/"
    type = "http"
}

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
