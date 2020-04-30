## Dynamic URL

Kori is working on a jq script she's going to use when hacking around with
some bash. She doesn't want to deal with authentication or risk getting rate
limited while she's working on it. So she uses mock-proxy to substitute for
the real thing.

Using the [CircleCI documentation](https://circleci.com/docs/api/v2/#get-recent-runs-of-a-workflow)
she generates a Routes file:

```hcl
route {
  host = "circleci.com"
  path = "/api/v2/insights/:project/workflows/:workflow"
  type = "http"
}
```

And a mock file, at `/mocks/cirleci.com/api/v2/insights/:project/workflows/:workflow.mock`.

The use of `:project` and `:workflow` in the route and mock file definitions
will automatically create substitutions for the `{{.project}}` and
`{{.workflow}}` values in the mock file when a request is made that matches the
pattern.

```json
{
  "items": [
    {
      "id": "{{.project}}/{{.workflow}}-1",
      "duration": 0,
      "created_at": "2020-04-29T16:35:21Z",
      "stopped_at": "2020-04-29T16:35:21Z",
      "credits-used": 0,
      "status": "success"
    },
    {
      "id": "{{.project}}/{{.workflow}}-2",
      "duration": 0,
      "created_at": "2020-04-29T16:35:21Z",
      "stopped_at": "2020-04-29T16:35:21Z",
      "credits-used": 0,
      "status": "running"
    } 
  ],
  "next_page_token": "string"
}
```

Now once she starts mock-proxy, she can iterate on her jq script without
needing to hit the real CircleCI endpoints.

```bash
export http_proxy=http://squid.proxy:8888
export https_proxy=http://squid.proxy:8888

curl --silent https://circleci.com/api/v2/insights/myproject/workflows/myworkflow | \
  jq -r '.items[] | "\(.id) \(.status)"'

### myproject/myworkflow-1 success
### myproject/myworkflow-2 running

curl --silent https://circleci.com/api/v2/insights/otherproject/workflows/otherworkflow | \
  jq -r '.items[] | "\(.id) \(.status)"'

### otherproject/otherworkflow-1 success
### otherproject/otherworkflow-2 running
```
