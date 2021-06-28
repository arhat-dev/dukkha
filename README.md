# dukkha `苦`

[![CI](https://github.com/arhat-dev/dukkha/workflows/CI/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/dukkha/workflows/Build/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/dukkha)](https://pkg.go.dev/arhat.dev/dukkha)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/dukkha)](https://goreportcard.com/report/arhat.dev/dukkha)
[![codecov](https://codecov.io/gh/arhat-dev/dukkha/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-dev/dukkha)

Make YAML files Makefiles

## The Idea

- Type checked configuration (e.g. workflow definition for github actions)
- Language or tool specific support (e.g. `goreleaser` for go/docker/npm builds)
- Flexible scripting (e.g. Makefile, shell scripts)

A typical build automation tool only take one or two from the above, but we'd take three in `dukkha`!

## Features

- Rendering suffix (`@<renderer>`), get configuration updated dynamically at runtime
- Customizable task matrix execution everywhere
- Shell completion for tools, tasks and task matrix

## Known Limitations

- Rendering suffix:
  - Cannot render the value when applied in yaml file to a map key (in source code), see [#22](https://github.com/arhat-dev/dukkha/issues/22)

## LICENSE

```text
Copyright 2020 The arhat.dev Authors.

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
