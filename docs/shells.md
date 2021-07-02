# Shells

```yaml
# shell command to provision a shell to execute your scripts
shells:
# fist shell will become the default shell, you can reference it using `shell`,
# but other shells with `shell:<shell_name>`
- name: bash
  # environment variables used when provisioning a shell
  env:
  - SHELL_ENV_NAME=value
  cmd:
  - bash
  - -x
  - -c
- name: python
  cmd@shell: |-
    echo "- $(command -v python3)"
```
