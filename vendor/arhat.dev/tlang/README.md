# tlang

A simple scripting/templating language derived from golang template.

- [syntax](./docs/syntax.md)

## Performance

benchmark metrics on m1 mac:

```text
BenchmarkParsing_text/parse-8     1297  927683 ns/op  317473 B/op   6025 allocs/op
BenchmarkParsing_text/exec-8      4029  293519 ns/op   32061 B/op   2001 allocs/op
BenchmarkParsing_tlang/parse-8    2654  449507 ns/op  315122 B/op   6018 allocs/op
BenchmarkParsing_tlang/exec-8     4110  288588 ns/op   32060 B/op   2001 allocs/op
```

The improvement of parsing performance probably comes from the synchornoized [lexer](./parse/lex.go).

see [benchmark source code](./benchmark/bench_test.go)

## Credits

This project won't be possible without the inspiration of the `text/template` package from golang standard library.

## LICENSE

MIT
