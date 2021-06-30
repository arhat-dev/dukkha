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

### Task `docker:build`

Build container images

```yaml
docker:build:
- name: example-image
  # images_names of the build output
  # if not set, will use the `name` value as `image`
  image_names:
  - image: example.com/image:tag-amd64
    manifest: example.com/image:manifest-tag

  # docker build [options] <the only positional-arg is the context>
  context: "."

  # arg to docker build -f
  dockerfile: path/to/dockerfile

  # extra docker build args
  extraArgs: []
```

### Task `docker:push`

Push images and manifests

```yaml
docker:push:
- name: foo
  image_names:
  - image: example.com/foo:latest-amd64
    manifest: example.com/foo:latest
  extraArgs: []
```
