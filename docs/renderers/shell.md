# Shell Renderer

```yaml
foo@shell: echo "Woo"
```

Run embedded shell script and use the output to stdout as the real value

## Config Options

```yaml
renderers:
  # no options
  shell: {}
```

## Supported value types

- `string`

## Variants

### Variant `@shell:<shell-name>`

Use specific shell to run the script

Example:

```yaml
shells:
- name: sh
- name: python

tool:task-kind:
- name: foo
  foo@shell:python: print("hello")
```

## Suggested Use Cases

No suggestion for now
