# Renderers

## Renderer Attributes

A renderer attribute set additional instructions to renderer

### Format

`<renderer>#<attr_1,attr_2,...>`

Examples:

- Renderer `my-rdr` with no attributes: `foo@my-rdr`
- Renderer `my-rdr` with empty attributes: `foo@my-rdr#` (same as no attributes)
- Renderer `my-rdr` with attributes `a`, `b`, `c`: `foo@my-rdr#a,b,c`

__NOTE:__ There can be white-spaces around attributes (e.g. `foo@my-red#a, b, c` is equivalent to `foo@my-rdr#a,b,c`).
