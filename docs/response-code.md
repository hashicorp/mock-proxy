## Mocking Different Response Codes

By default, all mocks return a 200 when they succeed. That's not the only
possible thing you might need to mock though. In order to request a different
response code along for a mock, send along the header
`X-Desired-Response-Code`.

```
./hack/local-dev-up.sh
curl --head --header "X-Desired-Response-Code: 204" example.com
```
