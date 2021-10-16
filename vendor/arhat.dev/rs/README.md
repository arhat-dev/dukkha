# rs

[![CI](https://github.com/arhat-dev/rs/workflows/CI/badge.svg)](https://github.com/arhat-dev/rs/actions?query=workflow%3ACI)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/rs)](https://pkg.go.dev/arhat.dev/rs)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/rs)](https://goreportcard.com/report/arhat.dev/rs)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_rs?branch=master&token=563ff8cf318c9303285b4dd4eeb0c660)](https://sonar.arhat.dev/dashboard?id=arhat-dev_rs)

`R`endering `S`uffix Support for yaml

## What?

Before we start, let's agree on the concept `renderer` is just a simple function that takes a input and generates a output.

```text
input -> [ renderer ] -> output
```

A plain yaml doc looks like this:

```yaml
foo: bar
```

everything's static, once unmarshaled, `foo` gets value `bar`.

Now what if we would like to take an environment variable as `foo`'s value?

```yaml
foo: ${FOO}
```

you will only get `${FOO}` after unmarshaling, you have to code you own logic to map `${FOO}` to some system environment variable.

What you code, is actually a `renderer` with its usage specific to `foo`.

Rendering suffix is there to help, it offers a way to make your yaml doc dynamic on its own, you control what you get in yaml rather than your compiled code.

### How it looks?

Let's continue with the environment variable support, say we have developed a renderer `env` (which calls `os.ExpandEnv` to map env reference to its value)

```yaml
foo@env: ${FOO}
```

Before I would start, probably you have already figured out what's going on:

`@` is like a notifier, and notifies the renderer `env` with value `${FOO}`, then `env` does its own job and generates some output with `os.ExpandEnv` function call, and at last, the output is set as `foo`'s value.

Simple and straightforward, right?

### Why?

Both writing and resolving software configuration are tedious and error prone, especially when you want to achieve felixibility to some extent.

Here is an example, given target config like this

```go
type Config struct {
    A string
    B int
    C float64

    // RemoteHostScript for remote host execution
    RemoteHostScript string `yaml:"remote_host_script"`
}
```

when you want to support environment variables for its fields `A`, `B`, `C`, while not including `remote_host_script` (obviously the `remote_host_script` should be executed in some remote system), what whould you do?

```yaml
a: ${ENV_FOR_A}
b: ${ENV_FOR_B}
c: ${ENV_FOR_C}

remote_host_script: |-
  #!/bin/sh

  echo ${ENV_WITH_DIFFERENT_CONTEXT}
```

- Make `remote_host_script` a file path reference (or add a new field like `remote_host_script_file`)?
  - Simple and effective, but now you have two source of remote host script, more fixed code logic added, more document to come for preferred options when both is set. End user MUST read your documentation in detail for such subtle issues (if you are good at documenting).
- Expand environment variables before unmarshaling?
  - What would you do with `${ENV_WITH_DIFFERENT_CONTEXT}`? Well, you can unmarshal the config into two config, one with environment variables expanded, another not, and merge them into one.
- Unmarshal yaml as `map[string]interface{}` first, then do custom handling for every field?
  - Now you have to work with types manually, tedious yet error prone job starts now.
- Create a new DSL, add some keywords...
  - We have already seen so many DSLs created just for configuration purpose, almost none of them really simplified the configuration management, and usually only useful for development not deployment.

As developers, what we actually need is let end user decide which field is resolved by what method, and we just control when to resolve which field.

Rendering suffix is applicable to every single yaml field, doing exactly what end user need, and it can resolve fields partialy with some strategy in code, also exactly what developers want.

## Prerequisites

- Use `gopkg.in/yaml.v3` for yaml marshaling/unmarshaling

## Features

- Field specific rendering: customizable in-place yaml data rendering
  - Add a suffix starts with `@`, followed by some renderer name, to your yaml field name as rendering suffix (e.g. `foo@<renderer-name>: bar`)
- Type hinting: keep you data as what it supposed to be
  - Add a type hint suffix `?<some-type>` to your renderer (e.g. `foo@<renderer-name>?[]obj: bar` suggests `foo` should be an array of objects using result from `<renderer-name>` generated with input `bar`)
  - [list of supported type hints](https://github.com/arhat-dev/rs/blob/v0.4.0/typehint.go#L24))
- Data merging and patching made esay: create patching spec in yaml doc
  - Add a patching suffix `!` to your renderer (after the type hint if any), feed it a [patch spec](https://pkg.go.dev/arhat.dev/rs#PatchSpec) object
- Renderer chaining: render you data with a rendering pipeline
  - join you renderers with pipes (`|`), get your data rendered through the pipeline (e.g. join three renderers `a`, `b`, `c` -> `a|b|c`)
- Supports arbitraty yaml doc without type definition in code
  - Use [`AnyObject`](https://pkg.go.dev/arhat.dev/rs#AnyObject) as `interface{}`
  - Use [`AnyObjectMap`](https://pkg.go.dev/arhat.dev/rs#AnyObjectMap) as `map[string]interface{}`
- Vanilla yaml, valid for all standard yaml parser
- Everything `gopkg.in/yaml.v3` supports are supported
  - Anchors, Alias, YAML Merges

Sample YAML Doc with all features above

```yaml
foo@a|b|c!: &foo
  value@http?[]obj: https://example.com
  merge:
  - value@file: ./value-a.yml
    select: |-
      [ .[] | sort ]
  patch:
  - { op: remove, path: /0/foo }
  - op: add
    path: /0/foo
    value: "bar"
    select: '. + "-suffix-from-jq-query"'
  select: |-
    { foo: .[0].foo, bar: .[1].bar }

bar@a|b|c: *foo
```

## Usage

See [example_test.go](./example_test.go)

__NOTE:__ You can find more examples in [`arhat-dev/dukkha`](https://github.com/arhat-dev/dukkha)

## Known Limitations

See [known_limitation_test.go](./known_limitation_test.go) for sample code and workaround.

- Built-in map data structure with rendering suffix applied to map key are treated as is, rendering suffix won't be recognized.
  - which means for `map[string]interface{}`, `foo@foo: bar` is just a map item with key `foo@foo`, value `bar`, no data to be resolved

## LICENSE

```txt
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
