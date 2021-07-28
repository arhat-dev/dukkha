# dukkha `苦`

[![CI](https://github.com/arhat-dev/dukkha/workflows/CI/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/dukkha/workflows/Build/badge.svg)](https://github.com/arhat-dev/dukkha/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/dukkha)](https://pkg.go.dev/arhat.dev/dukkha)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/dukkha)](https://goreportcard.com/report/arhat.dev/dukkha)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_dukkha?branch=master&token=1f8a3998312d6feee60ab16f1ef58ca8)](https://sonar.arhat.dev/dashboard?id=arhat-dev_dukkha)

Make YAML files Makefiles

## Goals

- Type checked configuration (e.g. workflow definition for github actions)
- Language and tool specific support (e.g. `goreleaser` for go/docker/npm builds)
- Flexible scripting (e.g. `Makefile`, shell scripts)

A typical build automation tool only take one or two from the above, but we'd take three in `dukkha`!

## Non-Goals

- To replace any existing tool
  - `dukkha` only wraps other tools for common use cases to ease your life with devops pipelines
- To build a custom cli tool to support all kinds of tasks

## Features

- Rendering suffix, get configuration updated dynamically at runtime
  - Just `@some-renderer` in your yaml field key
    - e.g. `foo@env: ${FOO}` will expand all referenced the environment variable
  - Also supports renderer chaning
    - e.g. `foo@http|template|env` first to fetch content from remote http server, then execute it as a go template, and finally expand environment variables in the resulted data
- Flexible yet strict typing, not always string value for renderer input
  - Some Renderer can accept any kind of value to keep your yaml file highlighted as it should be
- Customizable task matrix execution everywhere
- Shell completion for tools, tasks and task matrix

```yaml
workflow:run:
- name@env: ${WORKFLOW_NAME}
  matrix:
    kernel: [linux]
    arch: [amd64]
  jobs@template:
  - shell@env: |-
      echo ${MATRIX_KERNEL}/{{ .Env.MATRIX_ARCH }}
  - shell@template|http: |-
      https://gist.githubusercontent.com/arhatbot/{{- /* line join */ -}}
      d1f27e2b6d7e41a7c9d0a6ef7e39a921/raw/{{- /* line join */ -}}
      1e014333a3d78ac1139bc4cab9a68685e5080685/{{- /* line join */ -}}
      echo.sh
```

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
