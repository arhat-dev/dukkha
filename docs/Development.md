# Development

## `dukkha` struct field tag

The `dukkha` field tag is intended for config parsing

```go
type Foo struct {
    field.BaseField

    Bar []string `dukkha:"other"`
}
```
