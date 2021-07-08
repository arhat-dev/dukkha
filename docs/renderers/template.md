# Template Renderer

```yaml
foo@template: |-
  {{ .Env.MATRIX_ARCH }}
```

Execute an embedded [golang template](https://golang.org/pkg/text/template/) and use the result as the real value

## Config Options

```yaml
renderers:
  # no options
  template: {}
```

## Supported value types

Any valid yaml value

## Variants

None

## Suggested Use Cases

No suggestion for now
