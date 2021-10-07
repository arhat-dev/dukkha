# `buildah`

Build OCI images using [buildah](https://github.com/containers/buildah)

## Config

```yaml
tools:
  buildah:
  - name: <name your buildah tool>
    env: []
  - name: <another buildah tool>
    cmd:
    # an example to run buildah in docker
    - docker
    - run
    - -it
    - --rm
    - --workdir
    - $(pwd)
    - -v
    - $(pwd):$(pwd)
    - --security-opt
    - label=disable
    - --security-opt
    - seccomp=unconfined
    - -v
    - buildah-storage:/var/lib/containers
    - --device
    - /dev/fuse:rw
    - quay.io/buildah/stable
    - buildah
```

## Supported Tasks

### Task `buildah:login`

Login to registries

```yaml
buildah:login:
- name: example-login
  registry: ghcr.io
  username@env: ${GHCR_USER}
  # password is always passed to buildah via --password-stdin
  password@env: ${GHCR_PASS}
  # tls_skip_verify: false
```

### Task `buildah:build`

Build OCI images using `buildah bud` (bud: build-using-dockerfile)

```yaml
buildah:build:
- name: example-image
  # images_names of the build output
  # if not set, will use the `name` value as `image`
  image_names:
  - image: example.com/image:tag-amd64
    manifest: example.com/image:manifest-tag

  # if there is no tag set to the name (`:<some tag>` suffix), dukkha will set its tag
  # based on GIT_BRANCH, GIT_DEFAULT_BRANCH, GIT_TAG, GIT_WORKTREE_CLEAN,
  # GIT_COMMIT and MATRIX_ARCH, which we believe is suitable for most projects
  - image: defaulting-tag.example.com/image

  # buildah build [options] <the only positional-arg is the context>
  context: "."

  # arg to buildah build -f
  file: path/to/dockerfile

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

### Task `buildah:xbuild`

Build OCI images in yaml without dockerfile

```yaml
buildah:xbuild:
- name: example-image
  # images_names of the build output
  # if not set, will use the `name` value as `image`
  image_names:
  - image: example.com/image:tag-amd64
    manifest: example.com/image:manifest-tag

  # if there is no tag set to the name (`:<some tag>` suffix), dukkha will set its tag
  # based on GIT_BRANCH, GIT_DEFAULT_BRANCH, GIT_TAG, GIT_WORKTREE_CLEAN,
  # GIT_COMMIT and MATRIX_ARCH, which we believe is suitable for most projects
  - image: defaulting-tag.example.com/image

  # build steps
  steps:
  - # each step may have unique id set manually (like FROM ... AS <id> in dockerfile)
    id: foo
    # from statement [buildah-from](https://github.com/containers/buildah/blob/main/docs/buildah-from.1.md)
    from:
      # ref of image source, conform to formats supported by buildah from
      # https://github.com/containers/buildah/blob/main/docs/buildah-from.1.md#description
      # - `dir:path`
      # - `docker://docker-reference`
      # - `docker-archive:path`
      # - `docker-daemon:docker-reference`
      # - `oci:path:tag`
      # - `oci-archive:path:tag`
      #
      # or just plain old docker image reference
      #
      # - `registry.domain/image-name`
      ref: example.com/foo

      # kernel value to set --os flag
      kernel: ""
      # arch value to set --arch and --variant flag)
      arch: ""

      # extra buildah pull flags
      extra_pull_args: [] # string

      # extra buildah from flags
      extra_args: [] # strings
  - # name is just a human readable description of this step
    name: Do something
    # record this step into history (defaults to true)
    record: false
    # commit the container finishing this step
    commit: false
    # image name of this step if commit is set to true
    #
    # ignored when this step is the last step
    commit_as: some-intermediate-image:foo
    # extra buildah commit flags
    extra_commit_args: [] # string

    run:
      # run some script, MUST set shebang
      script: |-
        #!/bin/sh

        echo FOO > /foo

      # extra buildah run flags
      extra_args: [] # strings

  - run:
      # or just run some bare commands not using shell
      cmd:
      - ps

  - from:
      ref: scratch

  - copy:
      # copy from previous step
      from:
        step:
          id: foo
          path: /foo
      to:
        path: /foo

      # extra args for copy
      extra_args: [] # strings
  - copy:
      # copy from some image
      from:
        image:
          ref: docker.io/library/alpine
          path: /bin/apk

          # extra buildah pull flags
          extra_pull_args: [] # strings
      to:
        path: /bin/apk
  - copy:
      # copy from local path
      from:
        local:
          path: ./docs/tools/buildah.md
      to:
        path: /docs/buildah.md
  - copy:
      # copy from http endpoint
      from:
        http:
          url: https://example.com/some-file
      to:
        path: /some-http-data

    # set image config (runs buildah config)
  - set:
      # --workdir
      workdir: /docs
      # --user
      user: foo:bar
      # --shell
      shell: [] #strings
      # --env
      env:
      - name: FOO
        value: bar
      # --annotation
      annotations:
      - name: anno
        value: val
      # --label
      labels:
      - name: label
        value: val
      # --port
      ports: # strings
      - "8080"
      # --entrypoint
      entrypoint:
      - /usr/local/entrypoint.sh
      # --cmd
      cmd:
      - do
      - something
      volumes: # strings
      - /foo
      - /bar
      stop_signal: SIGINT
```
