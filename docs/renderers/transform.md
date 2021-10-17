# Transform Renderer

```yaml
foo@transform:
  value: foo
  ops:
  - template: |-
      {{ .Value }}-do-something
```

Use operations like [template](https://golang.org/pkg/text/template/) to transform string value into arbitrary valid yaml or string, and use the result as the field value

## Config Options

```yaml
renderers:
  # no options
  transform: {}
```

## Supported value types

Only supports valid transform spec yaml object

```yaml
# value is a string value
value: String Only, seriously
# operations you want to take on the value
ops:
# currently template is the only operation we support
- template: |-
    add some {{- /* go */ -}} template
    your value above is available as {{ .Value }}
```

## Variants

None

## Suggested Use Cases

Convert your non-yaml data to yaml right in yaml.
