# Input Renderer

```yaml
foo@input: "Please Enter Something: "
```

Read user console input as field value.

## Config Options

__NOTE:__ Configuration is required to activate this renderer.

```yaml
renderers:
  input:
    # hide user input (e.g. entering password)
    hide: true
```

## Supported value types

- String: The input prompt

```yaml
foo@input: Enter your mood
```

## Suggested Use Cases

Force user interaction before task execution.
