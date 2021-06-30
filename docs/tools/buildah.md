# `buildah`

Build OCI images using [buildah](https://github.com/containers/buildah)

## Config

```yaml
tools:
  buildah:
  - name: <name your default buildah tool>
    env: []
    cmd: []
    # - buildah
  - name: <another buildah tool>
    cmd: []
    # an example to run buildah in docker
    # - |-
    #   docker run -it --rm \
    #     --workdir $(pwd) \
    #     -v $(pwd):$(pwd) \
    #     --security-opt label=disable \
    #     --security-opt seccomp=unconfined \
    #     -v buildah-storage:/var/lib/containers \
    #     --device /dev/fuse:rw \
    #     quay.io/buildah/stable \
    #     buildah
```

## Supported Tasks

### Task `buildah:bud`

Build OCI images using `buildah bud` (bud: build-using-dockerfile)

Config is the same as [`docker:build`](./docker.md#task-dockerbuild), but replace `build` with `bud`, and `docker` with `buildah` in your mind

### Task `buildah:push`

Push images bulit with `buildah:bud` task to registries

Config is the same as [`docker:push`](./docker.md#task-dockerpush), but replace `docker` with `buildah` in your mind
