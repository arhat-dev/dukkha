# File Renderer

```yaml
foo@file: path/to/local/file
```

Read local file content as the real value

## Config Options

```yaml
renderers:
  file:
    # enable in memory cache, or always read local file
    enable_cache: true
    cache_max_age: "0"
    # in memory cache size limit
    cache_size_limit: 100M
```

## Supported value types

- `string`

## Suggested Use Cases

Config reuse
