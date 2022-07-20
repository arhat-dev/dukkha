# Tengo Script Renderer `tengo`

```yaml
foo@tengo: |-
  fmt := import("fmt")
  fmt.Println("hello")
```

Run [tengo](https://github.com/d5/tengo) script and use the output to stdout as the field value

## Config Options

```yaml
renderers:
  # no options
- tengo: {}
```

## Supported value types

- String: To run a single script
- List of Strings: To run a series of scripts in order

## Interoperation with `tmpl` renderer

You can call template funcs directly

```yaml
foo@tengo: archconv.DebianTripleName "armv6"
```

## Suggested Use Cases

TBD
