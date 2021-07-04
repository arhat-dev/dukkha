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
    # an example to run buildah in buildah
    # - |-
    #   buildah run -it --rm \
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

```yaml
buildah:bud:
- name: example-image
  # images_names of the build output
  # if not set, will use the `name` value as `image`
  image_names:
  - image: example.com/image:tag-amd64
    manifest: example.com/image:manifest-tag

  # if there is no tag set to the name (`:<some tag>` suffix), dukkha will set its tag
  # based on GIT_BRANCH, GIT_DEFAULT_BRANCH, GIT_TAG, GIT_WORKTREE_CLEAN
  # and GIT_COMMIT, which we believe is suitable for most projects
  - image: defaulting-tag.example.com/image

  # buildah build [options] <the only positional-arg is the context>
  context: "."

  # arg to buildah build -f
  dockerfile: path/to/dockerfile

  # extra buildah build args
  extra_args: []
```

### Task `buildah:push`

Push OCI images and manifests built by buildah to registries

```yaml
buildah:push:
- name: foo
  image_names:
  # only image/manifest names with FQDN as first part will be pushed
  - image: example.com/foo:latest-amd64
    manifest: example.com/foo:latest
```
