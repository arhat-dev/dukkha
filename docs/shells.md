# Shells

Define external shells to run scripts using renderer `shell:<shell-name>` or as action

```yaml
shells:
- name: bash
  # extra shell args
  env:
  - name: ENV_NAME
    value: env-value

  cmd:
  - bash
  - -x
  - -c

- name: python
  cmd@shell: |-
    echo "- $(command -v python3)"
```
