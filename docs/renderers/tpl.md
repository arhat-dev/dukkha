# Template Renderer `tpl`

```yaml
foo@tpl: |-
  {{ matrix.arch }}
```

Execute an embedded [golang template](https://golang.org/pkg/text/template/) and use the result as the real value.

__NOTICE:__

- See [template_funcs.md](../generated/template_funcs.md) for overview of supported functions
- Most template functions without namespace (i.e. functions without `.` in name, e.g. `deepEqual`) come from [`masterminds/sprig`](https://masterminds.github.io/sprig/)
- Most template functions with namespace (i.e. functions with `.` in name, e.g. `time.ParseLocal`) come from [`hairyhenderson/gomplate`](https://docs.gomplate.ca/)

## Config Options

```yaml
renderers:
  tpl:
    # include local template files, support glob pattern
    include:
    - foo/*.tpl
    - bar.tpl

    # a map of arbitrary values accessible from
    # template `{{ var.<key> }}`
    variables:
      # top level key must be string
      foo: [a, list, of, strings]
      bar: { an: { object: true } }
```

## Supported value types

- Any valid yaml value.
- Valid input spec (when `use-spec` attribute applied)

```yaml
foo@tpl#use-spec:
  # template to render
  template: |-
    {{- template "foo" -}}
    {{- var.data -}}

  # include local template files, support glob pattern
  #
  # NOTE: this option does not inherit renderer config include
  include:
  - foo/*.tpl
  - bar.tpl

  # a map of arbitrary values accessible from
  # template `{{ var.<key> }}`
  #
  # NOTE: this option does not inherit renderer config variable
  variables:
    # top level key must be string
    foo: [a, list, of, strings]
    bar: { an: { object: true } }
```

## Supported Attributes

- `use-spec`: Treat data to render as input spec instead of as template text.

## Interoperation with `shell` renderer

There is a template func `eval.Shell` for running shell commands in template.

```yaml
foo@tpl: |-
  {{- eval.Shell "echo 'Called From Template'" -}}
```

## Suggested Use Cases

No suggestion for now.
