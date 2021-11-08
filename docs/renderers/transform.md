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

Use operations like [template](https://golang.org/pkg/text/template/) to transform string value into arbitrary valid yaml or just plain string, and use the result as the field value.

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

- Convert your non-yaml data to yaml right in yaml.
- Composite different renderers to achieve significantly more flexibility.

  ```yaml
  # entrypoint is the transform renderer
  # last step happens here: archivefile renders value generated from renderer `transform`
  foo@transform|archivefile:
    # first step happens here: fetch data.tar.gz from remote http endpoint
    # notice the `#cached-file`, attribute `cached-file` will make renderer
    # `http` return local file path to the cached content.
    value@http#cached-file: https://example.com/data.tar.gz
    ops:
    # second step happens here: format the resolved `value` for render `archivefile`
    - template: |-
        {{ .Value }}:/in-archive-target-file
  ```
