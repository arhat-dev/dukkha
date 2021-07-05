# Shell Renderer

- `@shell`
- `@shell_file`

Run shell scripts and use the output to stdout as the real value

- [Supported value types](#supported-value-types)
- [Variants](#variants)
  - [Variant `@shell:<shell-name>`](#variant-shellshell-name)
- [Suggested Use Cases](#suggested-use-cases)

## Supported value types

Both `@shell` and `@shell_file`

- `string`

`@shell_file` only:

- `[]string`

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
