# File Renderer

```yaml
foo@file: path/to/local/file
```

Read local file content as the real value

## Config Options

```yaml
renderers:
- file:
    # enable in memory cache, or always read local file
    cache:
      enabled: true
      timeout: "0"
      size: 100M
```

## Supported value types

- String: Local file path

```yaml
foo@file: /tmp/data.json
```

## Suggested Use Cases

- Local config reuse
