# git

## Supported Tasks

### Task `git:clone`

Ensure the destination path is a repo cloned from remote and clean

```yaml
git:clone:
- name: example
  url: https://example.com/repo.git
  path: ./third_party/repo
  # remote branch to checkout, defaults to the remote default branch
  remoteBranch: v0.1.0
  # local branch name checked out from the remote branch
  localBranch: ""
  # value to git clone --origin
  remoteName: upstream
  # extraArgs for clone
  extraArgs:
  - --depth=1
```
