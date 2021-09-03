# Plugins

## Plugin Types

- Renderer Plugin
- Tool Plugin

## Install A Plugin

```yaml
renderer_plugins:
- name: foo-renderer
  module: example.com/foo-renderer
- name: bar-renderer
  source: |-
    bar-renderer

tool_plugins:
- name: foo-tool
  module: example.com/foo-tool
- name: bar-tool
  source: |-
    bar-tool
```
