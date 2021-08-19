# rs

[![CI](https://github.com/arhat-dev/rs/workflows/CI/badge.svg)](https://github.com/arhat-dev/rs/actions?query=workflow%3ACI)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/rs)](https://pkg.go.dev/arhat.dev/rs)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/rs)](https://goreportcard.com/report/arhat.dev/rs)
[![Coverage](https://badge.arhat.dev/sonar/coverage/arhat-dev_rs?branch=master&token=563ff8cf318c9303285b4dd4eeb0c660)](https://sonar.arhat.dev/dashboard?id=arhat-dev_rs)

`R`endering `S`uffix Support for yaml

## Usage

1. Embed `rs.BaseField` as the very first field in your struct

   ```go
   type MyStruct struct {
     rs.BaseField // MUST be the first one
   }
   ```

1. Write yaml doc with rendering suffix (say `my-renderer`)

   ```yaml
   foo@my-renderer: bar
   ```

1. Implement the `my-renderer` to comply interface `rs.RenderingHandler` requirments (say `MyRendererImpl`)

1. Unmarshal your yaml data using [`go-yaml`](https://github.com/go-yaml/yaml) to `MyStruct` object (say `myStructObj`)

1. Resolve your yaml data by calling `myStructObj.ResolveFields(&MyRendererImpl{}, -1)`

__NOTE:__ You can find more source code usage examples in [`arhat-dev/dukkha`](https://github.com/arhat-dev/dukkha)

## Known Limitations

- Built-in map data structure with rendering suffix applied to map key are not supported

  Sample Code:

  ```go
  type Bar struct {
    // ... any kind of fields
  }

  type Foo struct {
    rs.BaseField

    SomeBar map[string]Bar `yaml:"some_bar"`
  }
  ```

  ```yaml
  some_bar:
    # this will not be resolved
    some_key_for_bar@my-renderer: my-renderer-value
  ```

  Workaround: Define your own struct, embed `rs.BaseField` as the first field and make the map a field with `rs:"other"` tag

  ```go
  type BarMap struct{
    rs.BaseField

    Data map[string]Bar `rs:"other"`
  }

  type Foo struct {
    rs.BaseField

    SomeBar BarMap `yaml:"some_bar"`
  }
  ```

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
