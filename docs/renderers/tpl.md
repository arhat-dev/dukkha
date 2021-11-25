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
- tpl:
    # include local template files, with glob pattern support
    include:
    - path: foo/*.tpl
    - text@file?str: bar.tpl
    - text: |-
        {{- define "my-template" -}}
          sample data
        {{- end -}}

    # a map of arbitrary values accessible from
    # template `{{ var.<key> }}`
    variables:
      # top level key must be string
      foo: [a, list, of, strings]
      bar: { an: { object: true } }
```

__NOTE:__ Using `path` for template inclusion (`include` option) reads files from local filesystem every time it's rendering, to avoid this, use `text` with some cache enabled renderer.

## Supported value types

- Any valid yaml value.
- Valid input spec (when `use-spec` attribute applied)

```yaml
foo@tpl#use-spec:
  # template to render
  template: |-
    !! if you have {{ define }} block in your template,
    !! you can execute it by name
    {{- template "foo" -}}

    !! and access variables here and in included templates
    {{- var.data -}}

    !! included file templates without {{ define }} can also be executed by index or file basename
    {{- template "1" -}} will execute second file template
    {{- template "foo.tpl" -}} will execute last `foo.tpl`

    !! like file templates, included text templates without {{ define }} can be executed by index with prefix `#`
    {{- template "#1" -}} will execute second text template

  # include template files (with glob pattern support) and plain text tempaltes
  #
  # NOTE: this option does not inherit `include` option from renderer config
  include:
  - path: foo/*.tpl
  - text@file?str: bar.tpl
  - text: foo

  # a map of arbitrary values accessible from
  # template `{{ var.<key> }}`
  #
  # NOTE: this option does not inherit `variables` from renderer config
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
