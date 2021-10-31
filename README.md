# dukkha `苦`

[![CI](https://github.com/arhat-dev/dukkha/workflows/CI/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/dukkha/workflows/Build/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/dukkha)](https://pkg.go.dev/arhat.dev/dukkha)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/dukkha)](https://goreportcard.com/report/arhat.dev/dukkha)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_dukkha?branch=master&token=1f8a3998312d6feee60ab16f1ef58ca8)](https://sonar.arhat.dev/dashboard?id=arhat-dev_dukkha)

Make YAML files Makefiles

## Goals

- Type checked configuration (like workflow definitions for github actions)
- Language and tool specific support (like `goreleaser`'s support for go/docker/rpm)
- Flexible scripting (e.g. `Makefile`, shell scripts)
- Content generation and certral management of config recipes. (e.g. `jsonnet`, `cuelang`)

A typical build automation tool only takes one or two from the above at the same time, but we have them all in `dukkha` thanks to the [rendering suffix][rs] support.

## Features

### Content Rendering Features

- Rendering suffix
  - This is the way we make YAML files Makefiles, have a look at [arhat-dev/rs][rs] to familiar yourself with rendering suffix.
  - Renderers like `http`, `env`, `file`, `template` ... are available in dukkha as built-in renderers, see [docs/renderers](./docs/renderers) for more details.

- Available as a cli for custom content rendering, run `dukkha render` over your own yaml docs using rendering suffix.

- Editor support: Autocompletion for dukkha config (including autocompletion of patch spec)
  - Add `https://raw.githubusercontent.com/arhat-dev/dukkha/master/docs/generated/schema.json` to your yaml schemas
    - For vscode, add to `yaml.schemas` (for [yaml-language-server](https://github.com/redhat-developer/yaml-language-server#language-server-settings) which is embedded in [`redhat.vscode-yaml`](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) plugin)

### Task Execution Features

- Declarative task definition, context aware and customizable tool invocation.
  - Tools know context values of its task execution, and can adapt itself for special case handling like cross compiling.

- Embedded shell environment, you can completely forget external shells with dukkha
  - Predictable command execution made esay, no more worries about when will the environment variables get set, you gain total control.

- Customizable task matrix execution everywhere
  - While matrix values are still limited to strings in order to make them available as environment variables.
  - Command line option `--matrix` (`-m`) for `dukkha run` controls which matrix set is chosen.

  ```yaml
  workflow:run:
  - name: matrix-example
    matrix:
      # add your matrix spec
      kernel: [linux, windows]
      arch: [amd64, arm64]
      my-mat-entry: [foo, bar]

      # and exclude some
      exclude: # `exclude` is reserved
      # match certain matrix entry
      - { kernel: [windows], arch: [amd64] }
      # partial matching is supported as well
      - arch: arm64

      # and include extra matrix, regardless of what `exclude` says
      include: # `include` is reserved field
      - { kernel: [linux], arch: [x86, riscv64] }
      - { kernel: [darwin], arch: [arm64] }
  ```

- Shell completion for tools, tasks and task matrix
  - Run `dukkha completion --help` for instructions

- ANSI escape sequence handling for commands not respecting tty settings to avoid lengthy while meaningless log output (e.g. maven)
  - Automatically enabled when stdout/stderr is not a tty
  - Can be manually enabled by setting flag `--translate-ansi-stream` and `--retain-ansi-style` when running task
  - This functionality is largely based on [`github.com/aoldershaw/ansi`](https://github.com/aoldershaw/ansi)

## How tasks looks?

Here is just a `workflow` task

```yaml
workflow:run:
- name: quick-example
  matrix:
    kernel: [linux]
    arch: [amd64]
  jobs@template:
  # render environment variables before shell evaluation
  - shell@env: |-
      echo ${MATRIX_KERNEL}/{{ matrix.arch }}
  # run shell script from http server
  - shell@template|http: |-
      https://gist.githubusercontent.com/arhatbot/{{- /* line join */ -}}
      d1f27e2b6d7e41a7c9d0a6ef7e39a921/raw/{{- /* line join */ -}}
      1e014333a3d78ac1139bc4cab9a68685e5080685/{{- /* line join */ -}}
      echo.sh
```

__NOTE:__ You can find more [examples here](./docs/examples)

## Installation

### Build from source (not stripped)

- Option 1: Build directly with `go` 1.16+

  ```bash
  export VERSION=latest
  go get -u arhat.dev/dukkha/cmd/dukkha@${VERSION}
  ```

- Option 2: Clone and build with `make` and `go` 1.16+

  ```bash
  git clone https://github.com/arhat-dev/dukkha.git
  cd dukkha
  # checkout branch or commit as you prefer
  make dukkha
  # then you can find the built executable in build/dukkha
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
sget --key https://arhat.dev/.well-known/cosign.pub -o dukkha \
  "ghcr.io/arhat-dev/dist/dukkha:${VERSION}-${KERNEL}-${ARCH}"
chmod +x dukkha
```

- Option 2: Download pre-built artifacts from [releases](https://github.com/arhat-dev/dukkha/releases)

## Further Thoughts

As you may have noticed, the core of `dukkha` is the rendering suffix, we have more thoughts on the usage of rendering suffix as it dramatically eases the management of yaml config:

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

[rs]: https://github.com/arhat-dev/rs#readme
