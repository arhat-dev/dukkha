# golang

[Golang](https://golang.org/) toolchain support

`GO_COMPILER_PLATFORM="$(go version | cut -d\  -f4)"`

## Supported Tasks

### Task `golang:build`

Run go build

### Task `golang:test`

Run go test

### Task `golang:profile`

Run go tool pprof
