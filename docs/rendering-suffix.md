# Rendering Suffix

The way how `dukkha` make YAML files Makefiles

It's no magic, but a custom yaml unmarshaling method for all structs in `dukkha` with the help of [`field.BaseField`](https://pkg.go.dev/arhat.dev/dukkha/pkg/field#BaseField)

## What is a Rendering Suffix?

As indicated by the name, it's a suffix for rendering.

More specifically, it's a __name suffix__ in `@<renderer-name>` format to any yaml field name to set value generation engine (the renderer), but __WON't CHANGE the field name__ when parsing yaml docs.

Example:

Say you have a struct defined with `foo` field:

```go
type Example struct {
  Foo interface{} `yaml:"foo"`
}
```

In usual case, only yaml docs with fixed `foo: ...` can be parsed as value of `Example` type with `yaml.Unmarshal()`

But with rendering suffix support, yaml docs like `foo@my-renderer: ...` can also be parsed as the `Example` type.

## What Can Rendering Suffix Do?

To generate field value dynamically, aka. Conditional Rendering

## How Rendering Suffix Works?

For example:

Without Rendering Suffix, your yaml file is static, values are parsed as is, for example, parsing `foo: bar` and you will get foo=bar as a result in you application.

While with Rendering Suffix, values are evaluted in application with runtime context, so

```yaml
foo@shell: printf "bar"
```

- The suffix will not change the field name
- When using rendering suffix, `dukkha` will evaluate the environment variables after the execution.

## Supported Rendering Suffixes

### Env Renderer

Expand environment variables to generate real field value

- `@env`: value is expected to be a string containing environment vairable references
- `@env_file`: value is expected to be a local file path to file

### File Renderer

Read file content as the real field value

- `@file`: value is epxected to be a local file path

### Template Renderer

Evaluate golang templates as the real field value

- `@template`: value is expected to be go template string
- `@template_file`: value is expected to be a local file path of the golang template file

### Shell Renderer

Run shell scripts and use the output to stdout as the real field value

- `@shell` or `@shell:<shell_name>`: value is expected to be shell script
- `@shell_file` or `@shell_file:<shell_name>`: value is expected to be a local file path of the shell script
