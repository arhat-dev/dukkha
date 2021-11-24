# Env Renderer

```yaml
foo@env: some $ENV_NAME or ${REF} or $(shell command)
```

Generate field value by expanding environment variable references (`$env_name` or `${env_name}`) and evaluating shell command (`$(some command)`).

__NOTE:__ Backquoted string (e.g. <code>\`do something\`</code>) is not treated as shell evaluation, that's just plain text to this renderer.

## Config Options

```yaml
renderers:
- env:
    # enable arbitrary command execution (e.g. `$(do something)`) when expanding env
    #
    # Defaults to `false`
    enable_exec: true
```

## Supported value types

Any valid yaml value

## Suggested Use Cases

Gain fine grained control over environment variables' resolution, so you don't need external shell to execute command.

### Expand environment variables before task execution start

Case 1: In `buildah:build` tasks, manifest name is hashed to generate a local manifest name, if it includes environment variables not expanded, you may produce malformed manifest with unwanted images mixed.

To address this, use env renderer to expand environment variables:

```yaml
image_names@env:
- image: foo:$(git describe --tags -C /path/to/foo/src)-${MATRIX_ROOTFS}-${MATRIX_ARCH}
  manifest: foo:$(git describe --tags -C /path/to/foo/src)-${MATRIX_ROOTFS}
```

Case 2: In `buildah:login` tasks, password is passed to `buildah login` via stdin, you have to resolve the password value before running the task.

```yaml
password@env: ${MY_PASSWORD}
```
