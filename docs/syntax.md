# Syntax

- [Rendering Suffix `@<rendering_method>`](#rendering-suffix-rendering_method)

```yaml
foo: bar
```

`foo` is the field name, `bar` is the value

## Rendering Suffix `@<rendering_method>`

Every field name can have a `@xxx` suffix to imply its rendering method,

```yaml
foo@shell: printf "bar"
```

- The suffix will not change the field name
- When using rendering suffix, `dukkha` will evaluate the environment variables after the execution.

Supported suffixes are:

- `@template` will use go template for rendering
- `@template_url` is similar to `@template`, but will read the url target as the template
- `@shell` will execute the value as shell script and use the output to stdout as the final value, script execution failure will cause `dukkha` to abort the execution
  - `@shell:<shell_name>` is also supported
- `@shell_url` is similar to `@shell`, but will read the url target as the script
  - `@shell_url:<shell_name>` is supported as well
- `@file` will use local file content as value

for `@xxx_url`:

- `file://./path/to/some/file`
- `https://example.com/remote/file`

__NOTE:__ you have to solve character (un)escaping on your on own
