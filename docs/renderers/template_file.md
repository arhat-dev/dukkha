# Template File Renderer

```yaml
foo@template_file: path/to/local.tpl
```

Execute a [golang template](https://golang.org/pkg/text/template/) stored on local filesystem and use the result as the real value

## Config Options

```yaml
renderers:
  template_file:
    # enable in memory cache, or always read local file
    enable_cache: false
    #cache_max_age: 0
    # in memory cache size limit
    cache_size_limit: 100M
```

## Supported value types

- `string`

## Variants

None

## Suggested Use Cases

No suggestion for now
