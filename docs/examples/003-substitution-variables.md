## Substitution Variables

Perry is working on a system that parses HTML, and passes it through a series
of filter programs that may do additional work. He wants to test passing
special characters all the way through the network stack.

To this end Perry creates a Routes file:

```hcl
route {
  host = "example.com"
  path = "/"
  type = "http"
}
```

And a mock file, at `/mocks/example.com/index.mock`.

```
Hello, {{ .Name }}
```

Now he can iterate different possible characters by passing them as
"Substitution Variables".

```bash
curl -X POST -F "key=Name" -F "value=World" squid.proxy/substitution-variables
parse
### Hello, World
curl -X POST -F "key=Name" -F "value=>&2" squid.proxy/substitution-variables
parse
### Hello, >&2
curl -X POST -F "key=Name" -F 'value=!!_@3%2F' squid.proxy/substitution-variables
parse
### Hello, !!_@3%2F
```
