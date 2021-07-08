# HTTP Renderer

```yaml
foo@http: https://gist.github.com/my-gist
```

Render value using http GET

## Config Options

```yaml
renderers:
  http:
    # enable local cache, disable to always fetch from remote
    enable_cache: true
    cache_max_age: 1h
```

## Supported value types

- `string`

## Variants

None

## Suggested Use Cases

Config reuse
