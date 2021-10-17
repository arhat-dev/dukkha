# HTTP Renderer

```yaml
foo@http: https://example.com/some-file
```

Render value using http

## Config Options

__NOTE:__ Configuration is required to activate this renderer.

```yaml
renderers:
  http:
    # cache config
    # enable local cache, disable to always fetch from remote
    enable_cache: true
    cache_max_age: 1h

    # http config
    method: GET # if not set, defaults to GET
    user: basic-auth-username
    password: basic-auth-password
    headers:
    - name: User-Agent
      value: dukkha
    # body: ""
    tls:
      enabled: false
      ca_cert: |-
        <pem-encoded-ca-cert>
      cert: |-
        <pem-encoded-cert>
      key: |-
        <pem-encoded-cert-key>
      server_name: server-name-override
      # key_log_file: for-tls-debugging
      # cipher_suites: []
      # insecure_skip_verify: true
    proxy:
      enabled: false
      http: http://proxy
      https: https://proxy
      # no_proxy:
      cgi: false
```

## Supported value types

- `string` of target url
- yaml config

  ```yaml
  url: http://url.example
  config:
    # options are the same as Config Options .renderers.http
    # but without cache related options
    method: POST
  ```

## Suggested Use Cases

Organization to share build recipes using certral http service
