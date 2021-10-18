# Shell Renderer

```yaml
foo@shell: echo "Woo"
```

Run embedded shell script and use the output to stdout as the field value

## Config Options

```yaml
renderers:
  # no options
  shell: {}
```

## Supported value types

- String
- List of Strings

## Interoperation with `template` renderer

You can call template funcs by prefixing their names with `template:`

```yaml
foo@shell: template:archconv.DebianTripleName "armv6"
```

## Suggested Use Cases

No suggestion for now
