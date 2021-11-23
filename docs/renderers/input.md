# Input Renderer

```yaml
foo@input: "Please Enter Something: "
```

Read console input as value.

## Config Options

__NOTE:__ Configuration is required to activate this renderer.

```yaml
renderers:
  input:
    # hide user input (e.g. entering password)
    # defaults to false
    hide_input: true

    # input prompt message string
    prompt: ""
```

## Supported value types

- String: The input prompt

```yaml
foo@input: Enter your mood
```

- Valid Input Spec (only when applied with attribute `#use-spec`)

```yaml
foo@input#use-spec:
  # defaults to renderer config's hide_input option
  hide_input: true

  # defaults to renderer config's prompt option
  prompt: ""
```

## Suggested Use Cases

- Force user interaction before task execution.
- Interactive task automation.
