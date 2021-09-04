# Plugins

## Plugin Types

- Renderer Plugin: Add new renderer or override existing renderer
- Tool Plugin: Add new tool support
- Task Plugin: Extend or override existing tool tasks support

## Install A Plugin

```yaml
plugins:
  renderers:
  - name: foo-renderer
    module: example.com/foo-renderer
  - name: bar-renderer
    source: |-
      func NewRenderer_bar_renderer(name string) { ... }

  tools:
  - name: foo-tool
    tasks:
    - foo-task
    module: example.com/foo-tool
  - name: bar-tool
    tasks:
    - bar-task
    source: |-
      func NewTool_bar_tool() dukkha.Tool { ... }
      func NewTool_bar_tool_bar_task(name string) dukkha.Task { ... }

  tasks:
  - tool: foo-tool
    task: bar-task
    source: |-
      func NewTool_foo_tool_bar_task(name string) dukkha.Task { ... }
```
