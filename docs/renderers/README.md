# Renderers

Renderers are the core components to make dukkha config dynamic.

## Concepts

### Renderer Name

A renderer name consists of two parts:

- mandatory first part: kind name of the renderer
  - e.g. `af`, `tpl`, `http`
- optional second part: instance name of the renderer
  - e.g. `my-rdr-name`

If there is instance name, there MUST be a colon (`:`) between kind name and instance name.

Examples for renderer kind `tpl`:

- `tpl` (no instance name)
- `tpl:` (empty instance name, same as no instance name)
- `tpl:my-tpl` (with instance name `my-tpl`)

### Renderer Attributes

`<renderer>#<attr_1,attr_2,...>`

Renderer attributes are additional runtime rendering instructions, independent from renderer's own config and input data, intended to make renderers aware of what kind of input is, or what kind of output we expect.

Format: comma separated list of strings with start prefix `#`

Examples for renderer `tpl`:

- `tpl` (no attribute)
- `tpl#` (empty attributes, same as no attribute)
- `tpl#a,b,c` (attributes `a`, `b`, `c`)

## Configuration

Renderers are configured using top level `renderers` section in dukkha config:

```yaml
renderers: []
```

Each item of the `renderers` is a group of renderers not depending on each other:

Good:

```yaml
renderers:
# group #0
- env:foo: {}

# group #1
- # use env:foo defined in group #0
  file:my-file@env:foo: {}
```

Bad:

```yaml
renderers:
# group #0
- af:bar: {}
  # use af:bar in the same group to render
  # config of tpl can cause unexpected errors
  tpl:
    cache@af:bar: # ...
```
