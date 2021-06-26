# docker

Build docker images using [docker](https://github.com/containers/docker)

- [Config](#config)
- [Tasks](#tasks)
  - [Task `docker:build`](#task-dockerbuild)
  - [Task `docker:push`](#task-dockerpush)

## Config

```yaml
tools:
  docker:
  - name: <name your default docker tool>
    env: []
    # - DOCKER_BUILDKIT=1
    cmd: []
    # - docker
  - name: <another docker tool>
    cmd: []
    # - ssh
    # - remote-host
    # - docker
```
