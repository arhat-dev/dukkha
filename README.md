# dukkha `苦`

[![CI](https://github.com/arhat-dev/dukkha/workflows/CI/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/dukkha/workflows/Build/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/dukkha)](https://pkg.go.dev/arhat.dev/dukkha)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/dukkha)](https://goreportcard.com/report/arhat.dev/dukkha)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_dukkha?branch=master&token=1f8a3998312d6feee60ab16f1ef58ca8)](https://sonar.arhat.dev/dashboard?id=arhat-dev_dukkha)
[![Telegram](https://img.shields.io/static/v1?label=telegram&message=join&style=flat-square&logo=telegram&logoColor=ffffff&color=54A7E6&labelColor=555555)](https://t.me/joinchat/xcTR4nLDTOs2Yzcy)

Make YAML files Makefiles

## Goals

- Type checked configuration (like workflow definitions for github actions)
- Language and tool specific support (like `goreleaser`'s support for go/docker/rpm)
- Flexible scripting (e.g. `Makefile`, shell scripts)
- Content generation and certral management of config recipes. (e.g. `jsonnet`, `cuelang`)

A typical build automation tool only takes one or two from the above at the same time, but we have them all in `dukkha` thanks to the [rendering suffix][rs] support.

## Features

### Content Rendering Features

- Rendering suffix: extremely extensible & dynamic but type checked config wherever you want.
  - This is the way we make YAML files Makefiles, have a look at [arhat-dev/rs][rs] to familiar yourself with rendering suffix.
  - Renderers like `http`, `env`, `file`, `tpl` ... are available in dukkha as built-in renderers, see [docs/renderers](./docs/renderers) for more details.
  - In addition to basic renderer support, we have renderer attributes (`@<renderer>#<attr>`) to produce different kind of result
    - A common use case of this feature is to reuse renderer `http`: by default it returns the content fetched from remote endpoint, but when applied with attribute `cached-file` as `http#cached-file`, it will produce local file path to the cached content.

- Available as a cli for custom content rendering, run `dukkha render` over your own yaml docs using rendering suffix.

- Editor support: Autocompletion for dukkha config (including autocompletion of patch spec)
  - Add `https://raw.githubusercontent.com/arhat-dev/dukkha/master/docs/generated/schema.json` to your yaml schemas
    - For vscode, add to `yaml.schemas` (for [yaml-language-server](https://github.com/redhat-developer/yaml-language-server#language-server-settings) which is embedded in [`redhat.vscode-yaml`](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) plugin)

### Task Execution Features

- Declarative task definition, context aware tool invocation.
  - Tools have access to task related values, and can adapt itself for all use cases (e.g. cross compiling).

- Embedded cross-platform bash environment, you can completely forget external shells with dukkha
  - Predictable command execution made esay, gain full control on when and how the arguments get evaluated.

  __NOTE:__ You can still use external shells as long as you configure them in `shells` section.

- Matrix execution for every task
  - Use command line option `--matrix` (`-m`) to control which vectors are chosen.

  __NOTE:__ Vectors are currently limited to string values.

- Shell completion for defined tools, tasks and task matrix
  - Run `dukkha completion --help` for instructions

- ANSI escape sequence handling for commands not respecting tty settings to avoid lengthy while meaningless log output (e.g. maven)
  - Automatically enabled when stdout/stderr is not a tty
  - Can be manually enabled by setting flag `--translate-ansi-stream` and `--retain-ansi-style` when running task
  - This functionality is largely based on [`github.com/aoldershaw/ansi`](https://github.com/aoldershaw/ansi)

__NOTE:__ You can find more details in [docs](./docs)

## How it looks?

see [docs/examples](./docs/examples)

## Installation

### Build from source (not stripped)

Clone and build with `make` and `go` 1.17+

```bash
git clone https://github.com/arhat-dev/dukkha.git
cd dukkha
# checkout branch or commit as you prefer
make dukkha
```

Then you can find the built executable at `./build/dukkha`

### Download Pre-built Executables

- Option 1: Download and verify signature of dukkha using [`sget`][sget]

  ```bash
  # set dukkha version
  export VERSION=latest
  # set to your host kernel (same as GOOS value)
  export KERNEL=linux
  # set to your host arch (slightly different from GOARCH, see ./docs/constants.md)
  export ARCH=amd64

  sget --key https://arhat.dev/.well-known/cosign.pub -o dukkha \
    "ghcr.io/arhat-dev/dist/dukkha:${VERSION}-${KERNEL}-${ARCH}"
  chmod +x dukkha
  ```

  __NOTE:__ Please refer to [arhat-dev/dukkha-presets golang common matrix](https://github.com/arhat-dev/dukkha-presets/blob/dev/matrix/golang/1.17/common.yml) for `KERNEL` and `ARCH` values

- Option 2: Download signed artifacts from [releases](https://github.com/arhat-dev/dukkha/releases), then decompress the tarball/zipfile.

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
[sget]: https://github.com/sigstore/cosign#sget
