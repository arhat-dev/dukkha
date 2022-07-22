# Template Renderer `tmpl`

```yaml
foo@tmpl: |-
  {{ matrix.arch }}
```

Execute a [golang template](https://golang.org/pkg/text/template/) to generate value.

__NOTE:__ See [template_funcs.md](../generated/template_funcs.md) for overview of supported functions

## Config Options

```yaml
renderers:
- tmpl:
    # include local template files, with glob pattern support
    include:
    - path: foo/*.tmpl
    - text@file?str: bar.tmpl
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

__NOTE:__ When using `path` for template inclusion (`include` option) reads files from local filesystem every time it's rendering, to avoid this, use `text` with some cache enabled renderer.

## Supported value types

- Any valid yaml value.
- Valid input spec (when `use-spec` attribute applied)

```yaml
foo@tmpl#use-spec:
  # template to render
  template: |-
    !! if you have {{ define }} block in your template,
    !! you can invoke it by name
    {{- template "foo" -}}

    {{/* access variables using var */}}
    {{- var.data -}}

    {{/* included file templates can be invoked by index or file basename */}}
    {{- template "1" -}}         will invoke second file template
    {{- template "foo.tmpl" -}}  will invoke the last `foo.tmpl`

    {{/* included text templates can be invoked by index with additional prefix `#` */}}
    {{- template "#1" -}}        will invoke the second text template

  # include template files (with glob pattern support) and plain text tempaltes
  #
  # NOTE: this option does not inherit `include` option from renderer config
  include:
  - path: foo/*.tmpl
  - text@file?str: bar.tmpl
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
foo@tmpl: |-
  {{- (eval.Shell "echo 'Called From Template'").Stdout -}}
```

## Suggested Use Cases

TBD.
