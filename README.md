# dukkha `苦`

[![CI](https://github.com/arhat-dev/dukkha/workflows/CI/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/dukkha/workflows/Build/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/dukkha)](https://pkg.go.dev/arhat.dev/dukkha)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/dukkha)](https://goreportcard.com/report/arhat.dev/dukkha)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_dukkha?branch=master&token=1f8a3998312d6feee60ab16f1ef58ca8)](https://sonar.arhat.dev/dashboard?id=arhat-dev_dukkha)

Make YAML files Makefiles

## Quick Example

```yaml
workflow:run:
- name: quick-example
  matrix:
    kernel: [linux]
    arch: [amd64]
  # render inner go templates
  jobs@template:
  # render environment variables before shell evaluation
  - shell@env: |-
      echo ${MATRIX_KERNEL}/{{ .Env.MATRIX_ARCH }}
  # run shell script from http server
  - shell@template|http: |-
      https://gist.githubusercontent.com/arhatbot/{{- /* line join */ -}}
      d1f27e2b6d7e41a7c9d0a6ef7e39a921/raw/{{- /* line join */ -}}
      1e014333a3d78ac1139bc4cab9a68685e5080685/{{- /* line join */ -}}
      echo.sh
```

__NOTE:__ You can find more live examples [examples](./examples)

## Goals

- Type checked configuration (like workflow definitions for github actions)
- Language and tool specific support (like `goreleaser` support for go/docker/rpm builds)
- Flexible scripting in yaml files (e.g. `Makefile`, shell scripts)
- Content generation and certral management of config recipes

A typical build automation tool only take one or two from the above, but we'd take all in `dukkha`!

## Non-Goals

- To replace any existing tool
  - `dukkha` only wraps other tools for common use cases to ease your life with devops pipelines
- To build a custom cli tool to support all kinds of tasks

## Features

### Content Rendering Features

Please refer to [docs/rendering](./docs/rendering.md) for more details

- Rendering suffix: get configuration rendered dynamically at runtime
  - Just `@some-renderer` in your yaml field key

  ```yaml
  # environment variable FOO will be expanded and its value is used to set
  # `foo` field
  foo@env: ${FOO}
  ```

- Renderer chaning: Combine multiple renderers to achieve even more flexibility.
  - Chain your renderers with pipe symbol `|`

  ```yaml
  # in the following yaml example, dukkha will render the content three times
  #   1. (http) Fetch content from remote http server as specified by the url
  #   2. (template) Execute the fetched content as a go template
  #   3. (env) Expand environment variables in the resulted data
  foo@http|template|env: https://example.com/foo.yaml
  ```

- Patching and merging: Combine multiple yaml documents into one
  - Append suffix `!` to existing rendering suffix

  ```yaml
  foo@http!: # notice the suffix `!`
    # value for this renderer (http)
    value: https://example.com/bar.yaml
    # merge other yaml docs
    merge:
    - data@file: ./foo.yaml
    - data: [plain, yaml]
    # json patch (rfc6902) in yaml format
    patches:
    - { op: add, path: /a/b/c, value: foo }
  ```

- Available as a cli for custom content rendering
  - Run `dukkha render` over your own yaml docs using rendering suffix

### Task Execution Features

- Customizable task matrix execution everywhere

```yaml
workflow:run:
- name: matrix-example
  matrix:
    # add your matrix spec
    kernel: [linux, windows]
    arch: [amd64, arm64]

    # and exclude some
    exclude: # `exclude` is reserved
    # match certain matrix
    - kernel: [windows]
      arch: [amd64]
    # partial matching is supported as well
    - arch: arm64

    # and include extra matrix
    include: # `include` is reserved field
    - kernel: [linux]
      arch: [x86, riscv64]
    - kernel: [darwin]
      arch: [arm64]
```

- Shell completion for tools, tasks and task matrix
  - Run `dukkha completion --help` for instructions

## Installation

### Build from source

```bash
make dukkha
```

### Download Pre-built Executables

Before you start, set what version you would like to download

```bash
# set dukkha version
export VERSION=latest
# set to your host kernel (same as GOOS value)
export KERNEL=linux
# set to your host arch (slightly different from GOARCH, see ./docs/constants.md)
export ARCH=amd64
```

__NOTE:__ Combinations of `KERNEL` and `ARCH` are available at [scripts/dukkha/build-matrix.yml](./scripts/dukkha/build-matrix.yml)

- Option 1: Download and verify signature of dukkha using [`sget`](https://github.com/sigstore/cosign)

```bash
sget -key https://arhat.dev/.well-known/cosign.pub -o dukkha \
  "ghcr.io/arhat-dev/dist/dukkha:${VERSION}-${KERNEL}-${ARCH}"
chmod +x dukkha
```

## Further Thoughts

As you may have noticed, the core of `dukkha` is the rendering suffix, we have more toughts on the usage of rendering suffix as it dramatically eases the management of yaml config:

- Extend dukkha to be something like `systemd` for system configration and process management
- Manage kubernetes manifests in GitOps pipeline
- Infrastructure as Config instead of Infrastructure as Code

More to come up with ...

## LICENSE

```text
Copyright 2021 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
