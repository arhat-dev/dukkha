# Env Renderer

```yaml
foo@env: some $ENV_NAME or ${REF} or $(shell command)
```

Expand environment variable references (`$env_name` or `${env_name}`) and evaluate shell command (`$(some command)`) to generate real value

## Config Options

```yaml
renderers:
  # no options
  env: {}
```

## Suggested Use Cases

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
