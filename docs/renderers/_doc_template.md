# {{ .Name | title }} Renderer

- `@{{ .Name }}`
- `@{{ .Name }}:<varient-value>`

{{ .Description }}

- [Supported value types](#supported-value-types)
- [Variants](#variants)
  - [Variant `@{{ .Name }}:<variant-value>`](#variant--name-variant-value)
- [Suggested Use Cases](#suggested-use-cases)

## Supported value types

- `string`
- `[]string`

or

Any valid yaml value

## Variants

### Variant `@{{ .Name }}:<variant-value>`

Doc to the variant

Example:

```yaml
foo@{{ .Name }}:
```

## Suggested Use Cases
