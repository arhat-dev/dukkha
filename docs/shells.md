# Shells

```yaml
# shell command to provision a shell to execute your scripts
shells:
# reference your shells using `shell:<shell_name>` as renderer
# or hook action or job action (in workflow run)
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
