# Arbitrary File Rendering

Render any doc using rendering suffix

## Run

```bash
dukkha render ./source -r
```

## Example output

__NOTE:__ Depending on your host environment, some value may be different

```yaml
<!DOCTYPE html>
<html>
<body>
  <p>Example a</p>
  <p>Example b</p>
  <p>Example c</p>
  <p>Example d</p>
</body>
</html>
---
# Example Markdown Doc Generated In Yaml

## Go Template

Some value using go template `powerpc64le-linux-gnu`

## Environment Variables

Some value expanded from environment variable `master`
---
#!/bin/sh

set -eux

# shell evaluation in env renderer is disabled by default
# we should get what we write as is
foo="$(printf "%s" something)"

export foo
export current_host_kernel="darwin"
export current_host_arch="arm64"

```
