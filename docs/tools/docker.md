# docker

Build docker images using [docker (moby)](https://github.com/moby/moby)

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

## Supported Tasks

### Task `docker:login`

Login to registries

Config is the same as [`buildah:login`](./buildah.md#task-buildahlogin), but replace `buildah` with `docker` in your mind

__NOTE:__ docker login has no `skip_tls_verify` support

### Task `docker:build`

Build docker images

Config is the same as [`buildah:bud`](./buildah.md#task-buildahbud), but replace `bud` with `build`, and `buildah` with `docker` in your mind

### Task `docker:push`

Push docker images and manifests

Config is the same as [`buildah:push`](./buildah.md#task-buildahpush), but replace `buildah` with `docker` in your mind
