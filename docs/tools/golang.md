# golang

[Golang](https://golang.org/) toolchain support

`GO_COMPILER_PLATFORM="$(go version | cut -d\  -f4)"`

## Supported Tasks

### Task `golang:build`

Run go build

```yaml
golang:build:
- name: foo
  path: ./cmd/foo
  outputs:
  - build/foo
  cgo:
    enabled: false
    cflags: []
    cxxflags: []
    ldflags: [-static]
  # go -ldflags
  ldflags:
  - -s -w
  - -X "main.Version=v0.1.1"
```

### Task `golang:test`

Run go test

```yaml
golang:test:
- name: foo
  matrix:
    pkg@shell: find ./pkg -type d -
  package@env: ${MATRIX_PKG}
  cgo:
    enabled: false
    cflags: []
    cxxflags: []
    fflags: []
    ldflags: []
  ldflags: []
  race: true
  count: 1
  cpu: [1, 2, 4]
  parallel: 3
  failfast: false
  short: false
  timeout: 10m
  # match to run only matched tests
  match: ^Test.*$
  benchmark:
    enabled: false
    duration: 1h30s
    # match to run only matched benchmarks, if not set and
    # enabled is true, will run all benchmarks
    match: ^Benchmark.*$
  profile:
    coverage:
      enabled: true
      output: ./coverage.txt
      mode: atomic
      packages@env:
      - ${MATRIX_PKG}
  custom_args:
  - -foo
```

### Task `golang:profile`

Run go tool pprof
