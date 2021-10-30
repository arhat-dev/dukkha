# Template Renderer

```yaml
foo@template: |-
  {{ matrix.arch }}
```

Execute an embedded [golang template](https://golang.org/pkg/text/template/) and use the result as the real value.

__NOTICE:__

- [template_funcs.md](./template_funcs.md) provides an overview of supported functions
- Most template functions without namespace (i.e. functions without `.` in name, e.g. `deepEqual`) come from [`masterminds/sprig`](https://masterminds.github.io/sprig/)
- Most template functions with namespace (i.e. functions with `.` in name, e.g. `time.ParseLocal`) come from [`hairyhenderson/gomplate`](https://docs.gomplate.ca/)

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
