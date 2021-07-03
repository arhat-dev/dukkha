# github

[github-cli](https://github.com/cli/cli)

## Supported Tasks

### Task `github:release`

```yaml
github-cli:release:
- name: example
  tag: ${GIT_TAG}
  draft: true
  pre_release: true
  title: ""
  notes: ""
  files:
  - changelog.txt
```
