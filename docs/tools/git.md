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
  remote_branch: v0.1.0
  # local branch name checked out from the remote branch
  local_branch: ""
  # value to git clone --origin
  remote_name: upstream
  # extra args for git clone
  extra_args:
  - --depth=1
```
