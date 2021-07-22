# Rendering Suffix

The way how `dukkha` make YAML files Makefiles

It's no magic, but a custom yaml unmarshaling method for all structs in `dukkha` with the help of [`field.BaseField`](https://pkg.go.dev/arhat.dev/dukkha/pkg/field#BaseField)

## What Renderer Suffix is Not?

It's not a yaml extension, do not mess with the reserved yaml keyword `@`

## What is a Rendering Suffix?

It's a yaml parsing extension supported by `dukkha`.

As indicated by the name, it's a suffix used for rendering.

More specifically, it is:

- A __name suffix__ in `@<renderer-name>` format
  - `<renderer-name>` is just a placeholder
- Can be applied to any yaml field name
  - e.g. `foo@<renderer-name>: ...`
- To set value generation engine (the renderer) in dukkha
- But __WON't CHANGE__ the field name.

Example:

Say you have a struct defined with `foo` field:

```go
type Example struct {
  Foo interface{} `yaml:"foo"`
}
```

In usual case, only parsing yaml docs with exact `foo` key like

```yaml
foo: something something
```

to `Example` type using `yaml.Unmarshal()`(from [go-yaml](https://github.com/go-yaml/yaml)) can get the value of `foo`, any change to the field name will be treated as a different field

But with rendering suffix support, yaml docs like

```yaml
foo@my-renderer: woo
```

or

```yaml
foo@another-renderer: cool
```

(change the text between `@` and `:`, you get your own examples)

Can also be parsed as the `Example` type with `foo` field resolved in dukkha.

## What Can Rendering Suffix Do?

To generate field value dynamically, aka. Conditional Rendering

## How Rendering Suffix Works?

Without Rendering Suffix, your yaml file is static, values are parsed as is, for example, parsing `foo: bar` and you will get foo=bar as a result in you application.

While with Rendering Suffix, values are evaluted in application with runtime context, it's highly dynamic but also get type checked.

A simple example using embedded shell rendering suffix (assume a posix `sh`):

```yaml
foo@shell: printf "bar"
```

is equivalent to the yaml without rendering suffix:

```yaml
foo: bar
```
