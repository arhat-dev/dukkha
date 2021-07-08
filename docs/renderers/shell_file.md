# Shell File Renderer

```yaml
foo@shell_file: path/to/local/file
```

Run shell script stored on local filesystem and use the output to stdout as the real value

## Config Options

```yaml
renderers:
  shell_file:
    # enable in memory cache, or always read local file
    enable_cache: true
    #cache_max_age: 0
    # in memory cache size limit
    cache_size_limit: 100M
```

## Supported value types

- `string`
- array of `string`

## Variants

### Variant `@shell_file:<shell-name>`

Use specific shell to run the script

Example:

```yaml
shells:
- name: sh
- name: python
  cmd:
  - |-
    pipenv run python

tool:task-kind:
- name: foo
  foo@shell_file:python: ./some.py
```

## Suggested Use Cases

No suggestion for now
