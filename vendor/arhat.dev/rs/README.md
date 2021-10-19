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

Everything's static, once unmarshaled, `foo` gets value `bar`.

Now what if we would like to take an environment variable as `foo`'s value?

```yaml
foo: ${FOO}
```

you will only get `${FOO}` after unmarshaling, you have to code you own logic to map `${FOO}` to some system environment variable.

What you code, is actually a `renderer` with its usage specific to `foo`.

Rendering suffix is here to help, it offers a way to make your yaml doc dynamic on its own, you control what you get in yaml rather than your compiled code.

### How it looks?

Let's continue with the environment variable support, say we have developed a renderer `env` (which calls `os.ExpandEnv` to map env reference to its value)

```yaml
foo@env: ${FOO}
```

Before I would start, probably you have already figured out what's going on:

`@` is like a notifier, and notifies the renderer `env` with value `${FOO}`, then `env` does its own job and generates some output with `os.ExpandEnv` function call, and at last, the output is set as `foo`'s value.

But also wondering how can `foo@env` be resolved as `foo` since they are completely different field name! That's what we are taking care of, see [How it works?](#how-it-works) section for brief introduction.

It's simple and straightforward, right?

## Prerequisites

- Use `gopkg.in/yaml.v3` for yaml marshaling/unmarshaling

## Features

- Field specific rendering: customizable in-place yaml data rendering
  - Add a suffix starts with `@`, followed by some renderer name, to your yaml field name as rendering suffix (e.g. `foo@<renderer-name>: bar`)
- Type hinting: keep you data as what it supposed to be
  - Add a type hint suffix `?<some-type>` to your renderer (e.g. `foo@<renderer-name>?[]obj: bar` suggests `foo` should be an array of objects using result from `<renderer-name>` generated with input `bar`)
  - See [list of supported type hints](https://github.com/arhat-dev/rs/blob/v0.4.0/typehint.go#L24)
- Data merging and patching made esay: create patching spec in yaml doc
  - Add a patching suffix `!` to your renderer (after the type hint if any), feed it a [patch spec](https://pkg.go.dev/arhat.dev/rs#PatchSpec) object
- Renderer chaining: render you data with a rendering pipeline
  - join you renderers with pipes (`|`), get your data rendered through the pipeline (e.g. join three renderers `a`, `b`, `c` -> `a|b|c`)
- Supports arbitraty yaml doc without type definition in your own code.
  - Use [`AnyObject`](https://pkg.go.dev/arhat.dev/rs#AnyObject) as `interface{}`.
  - Use [`AnyObjectMap`](https://pkg.go.dev/arhat.dev/rs#AnyObjectMap) as `map[string]interface{}`.
- Everything `gopkg.in/yaml.v3` supports are supported.
  - Anchors, Alias, YAML Merges ...
- Extended but still vanilla yaml, your yaml doc stays valid for all standard yaml parser.

Sample YAML Doc with all features above

```yaml
foo@a!: &foo
  value@env|http?[]obj: https://example.com/${FOO_FILE_PATH}
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

bar@a!: *foo
```

__NOTE:__ This module provides no renderer implementation, and the only built-in renderer is a pseudo renderer with empty name that skips rendering (output is what input is) for data patching and type hinting purpose (e.g. `foo@?int!: { ... patch spec ... }`). You have to roll out your own renderers.

## Usage

See [example_test.go](./example_test.go)

__NOTE:__ You can find more examples in [`arhat-dev/dukkha`](https://github.com/arhat-dev/dukkha)

## Known Limitations

See [known_limitation_test.go](./known_limitation_test.go) for sample code and workaround.

- Golang built-in map with rendering suffix applied to map key are treated as is, rendering suffix won't be recognized.
  - For `map[string]interface{}`, `foo@foo: bar` is just a map item with key `foo@foo`, value `bar`, no data to be resolved.
  - The reason for this limitation is obvious since built-in map types doesn't have `BaseField` embedded, but it can be counterintuitive when you have a map field in a struct having `BaseField` embedded.

## How it works?

IoC (Inversion of Control) is famous for its application in DI (Dependency Injection), but what it literally says is to control the outside world at somewhere inside the world, that's the start point.

Custom yaml unmarshaling requires custom implementation of `yaml.Unmarshaler`, so usually you just get your structs unmarshaled by `yaml.Unmarshal` directly.

We implemented something called `BaseField`, it lives in your struct as a embedded field, all its methods are exposed to the outside world by default, and guess what, it implements `yaml.Unmarshaler`, so your struct implements `yaml.Unmarshaler` as well.

But can you control sibling fields in a struct? Not possible in golang unless with the help of outside world, that's why we need `Init()` function, calling `Init()` with your struct actually activates the `BaseField` in it, `Init()` function tells the inner `BaseField` what fields the parent struct have (its sibling fields), with the help of reflection.

You only have to call `Init()` once for the top level struct, since then the `BaseField` in it knows what to do with its sibling fields, it will also search for all structs with `BaseField` embedded when unmarshaling, call `Init()` for them, until the whole yaml doc is unmarshaled.

During the unmarshaling process, `BaseField.UnmarshalYAML` get called by `yaml.Unmarhsal`, it checks the input yaml field names, if a yaml field name has a suffix starting with `@`, then that yaml field will be treated as using rendering suffix, `BaseField` parses the yaml field name to know the real field name is (e.g. `foo@bar`'s real field name is `foo`) and sets the rendering pipeline with the suffix, it also saves the yaml field value on its own but not setting the actual strcut field, when you call `my_struct_with_BaseField.ResolveFields()`, it feeds the rendering pipeline with saved field value to generate actual field value and set that as struct field value.

All in all, `BaseField` handles everything related to yaml unmarshaling to support rendering suffix after initial activation, so all you need to do is to embed a `BaseField` as the very first field in your struct where you want to support rendering suffix and activate the top level struct (with `BaseField` embedded, which can be some inner field) with a `Init()` function call.

## FAQ

Have a look at [FAQ.md](./FAQ.md) or start/join a [discussion on github](https://github.com/arhat-dev/rs/discussions).

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
