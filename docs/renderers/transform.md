# Transform Renderer

```yaml
foo@transform:
  value: foo
  ops:
  - template: |-
      {{ .Value }}-do-something
  - shell: |-
      echo ${VALUE}
```

Use operations like [template](https://golang.org/pkg/text/template/) to transform string value into arbitrary valid yaml or string, and use the result as the field value

## Config Options

```yaml
renderers:
  # no options
  transform: {}
```

## Supported value types

- Valid transform spec in yaml

```yaml
foo@transform:
  # value is a string value
  value: String Only, seriously
  # operations you want to take on the value
  ops:
  # Execute golang template over .Value
  - template: |-
      add some {{- /* go */ -}} template
      your value above is available as {{ .Value }}
  # Execute shell script with env ${VALUE}
  - shell: |-
      echo "${VALUE}"
```

## Suggested Use Cases

Convert your non-yaml data to yaml right in yaml.
