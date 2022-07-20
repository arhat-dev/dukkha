# Shell Renderer

```yaml
foo@shell: echo "Woo"
```

Run bash script and use the output to stdout as the field value

## Config Options

```yaml
renderers:
  # no options
- shell: {}
```

## Supported value types

- String: To run a single script
- List of Strings: To run a series of scripts in order

## Interoperation with `tpl` renderer

You can call template funcs by prefixing their names with `tpl:`

```yaml
foo@shell: tpl:archconv.DebianTripleName "armv6"
```

## Suggested Use Cases

No suggestion for now.
