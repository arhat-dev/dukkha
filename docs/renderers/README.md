# Renderers

Renderers are the core components to make dukkha config dynamic.

## Concepts

### Name

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

### Alias

A renderer can have one alias of its renderer name

Some renderers work online (e.g. `http`, `git`), some work locally (e.g. `file`, `env`), alias makes the switching from one to another seamless.

Example config:

```yaml
renderers:
- http:
    alias: my-http
```

### Attributes

`<renderer>#<attr_1,attr_2,...>`

Renderer attributes are additional runtime rendering instructions, independent from renderer's own config and input data, intended to make renderers aware of what kind of input is, or what kind of output we expect.

Format: comma separated list of strings with start prefix `#`

Examples for renderer `tpl`:

- `tpl` (no attribute)
- `tpl#` (empty attributes, same as no attribute)
- `tpl#a,b,c` (attributes `a`, `b`, `c`)

### Default Attributes

You can set default attributes to renderer in renderer config, these default attributes are respected as long as there is no explicit attributes to the renderer.

Example:

Set default attributes

```yaml
renderers:
- http:
    attributes:
    - cached-file
```

Use the renderer with default attributes

```yaml
foo@http: https://example.com/content
# foo will be set to the path to local cache file of fetched content
```

Use the renderer with custom attributes (and no default attributes applied in this case)

```yaml
foo@http#cached-file,use-spec:
  url: https://example.com/content
```

__NOTE:__ to clear default attributes, use single `#` (e.g. `foo@http#`, which means reset attributes to `http` renderer)

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
