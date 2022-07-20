# Template Language Renderer `tl`

```yaml
foo@tl: |-
  define x
    val
  end
```

`tl` is a derived scripting language from [golang template](https://golang.org/pkg/text/template/), without the emphasis of __text templating__.

__NOTE:__ This renderer shares every feature of `tmpl` renderer, see [`tmpl`](./tmpl.md) for general idea of golang template.

## Config Options

```yaml
renderers:
- tl:
    # include local template files, with glob pattern support
    include:
    - path: foo/*.tl
    - text@file?str: bar.tl
    - text: |-
        define "my-template"
          sample data
        end

    # a map of arbitrary values accessible from
    # template `var.<key>`
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
foo@tl#use-spec:
  # template to render
  template: |-
    !! if you have {{ define }} block in your template,
    !! you can execute it by name
    template "foo"

    !! and access variables here and in included templates
    var.data

    !! included file templates without {{ define }} can also be executed by index or file basename
    template "1"
    !! will execute second file template
    template "foo.tl"
    !! will execute last `foo.tl`

    !! like file templates, included text templates without {{ define }} can be executed by index with prefix `#`
    template "#1"
    !! will execute second text template

  # include template files (with glob pattern support) and plain text tempaltes
  #
  # NOTE: this option does not inherit `include` option from renderer config
  include:
  - path: foo/*.tl
  - text@file?str: bar.tl
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
foo@tl: (eval.Shell "echo 'Called From Template'").Stdout
```

## Suggested Use Cases

TBD.
