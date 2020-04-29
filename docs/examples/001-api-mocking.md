## API Mocking

Aaron is working on a Python script that will hit an API endpoint his coworker
is still iterating on. His coworker sends over some mocks with the API response
structure.

Using this he generates a Routes file:

```hcl
route {
  host = "example.com"
  path = "/"
  type = "http"
}
```

And a mock file, at `/mocks/example.com/index.mock`.

```json
{
  "field1": 0,
  "field2": "string"
}
```

Now he can work on his script before the site is up, while still using the full
network stack of his program.

```python
import requests

proxies = {
    'http': 'http://squid.proxy:8888',
    'https': 'http://squid.proxy:8888',
}

resp = requests.get(
    'https://example.com',
    proxies=proxies,
    verify='/usr/local/share/ca-certificates/ca.pem'
)
print(resp.json())
```

```bash
python test.py
### {'field1': 0, 'field2': 'string'}
```
