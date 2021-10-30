# Template Renderer

```yaml
foo@template: |-
  {{ matrix.arch }}
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

## Interoperation with `shell` renderer

There is a template func `eval.Shell` for running shell commands in template

```yaml
foo@template: |-
  {{- eval.Shell "echo 'Called From Template'" -}}
```

## Suggested Use Cases

No suggestion for now.
