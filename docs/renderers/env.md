# Env Renderer

Expand environment variable references (`$env_name` or `${env_name}`) and evaluate shell command (`$(some command)`) to generate real value

- [Usage](#usage)
- [Supported value types](#supported-value-types)
- [Variants](#variants)
  - [Variant `@env:<shell-name>`](#variant-envshell-name)
- [Suggested Use Cases](#suggested-use-cases)
  - [Expand environment variables before task execution start](#expand-environment-variables-before-task-execution-start)

## Usage

- `@env`

## Supported value types

Any valid yaml value

## Variants

### Variant `@env:<shell-name>`

Depending on the `shells` section of your dukkha config, you can use shell specific env renderers with `@env:<shell-name>`, shell commands (`$(some command)`) will be handled using that shell.

Example:

```yaml
shells:
- name: bash
- name: python

tool:task-kind:
- name: foo
  foo@env: command will be handled by the default shell `bash` $(echo "bash")
  bar@env:bash: $(echo "Use bash!")
  foo_bar@env:python:
    a: command get rendered by python
    b: so you can run python commands $(print('hello'))
```

## Suggested Use Cases

### Expand environment variables before task execution start

Case 1: In `buildah:bud` tasks, manifest name is hashed to generate a local manifest name, if it includes environment variables not expanded, you may produce malformed manifest with unwanted images mixed.

To fix this, use env renderer to expand environment variables:

```yaml
image_names@env:
- image: foo:$(git describe --tags -C /path/to/foo/src)-${MATRIX_ROOTFS}-${MATRIX_ARCH}
  manifest: foo:$(git describe --tags -C /path/to/foo/src)-${MATRIX_ROOTFS}
```

Case 2: In `buildah:login` tasks, password is passed to `buildah login` via stdin, you have to resolve the password value before running the task.

```yaml
password@env: ${MY_PASSWORD}
```
