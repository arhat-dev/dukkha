# github

[github-cli](https://github.com/cli/cli)

## Supported Tasks

### Task `github:release`

Create release and upload assets

```yaml
github:release:
- name: example
  tag: ${GIT_TAG}
  draft: true
  pre_release: true
  title: ""
  notes: ""
  # upload files
  files:
    # path to the file, glob is supported
  - path: changelog.txt
    # display label
    label: CHANGELOG
  - path: build/*
    # if multiple file matches the glob, label will get indexed suffix
    # e.g. `build-asset 1`
    label: build-asset
```

__NOTE:__ This task will fail if executed multiple times (hint: avoid matrix usage)
