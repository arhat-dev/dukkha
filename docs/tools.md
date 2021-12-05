# Tools

## Config Declaration

Tools are configured in top-level `tools` config section

```yaml
tools:
  <tool-kind>: []
```

Example:

```yaml
tools:
  workflow: []
```

## Common Tool Options

- `name: string`: tool name that can be referenced in cli or task reference
- `env: []Env`: tool specific env
- `cmd: []string`: exec strings to run this tool (no env expansion)

Example:

```yaml
tools:
  golang:
  - name: local
    env:
    - name: GOSUMDB
      value: "off"
    cmd: [go]
```
