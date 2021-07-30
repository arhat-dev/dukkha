# cosign

Container Signing using [`cosign`](https://github.com/sigstore/cosign)

## Config

```yaml
tools:
  cosign:
  - name: <name your cosign tool>
    env: []
    # - name: COSIGN_EXPERIMENTAL
    #   value: "1"
  - name: <another cosign tool>
    # cmd: []
```

## Supported Tasks

### Task `cosign:upload`

Upload blob/wasm to OCI registry

```yaml
cosign:upload:
- name: foo
  # type of the uploaded content, blob or wasm
  kind: blob
  matrix:
    # this task recognizes `kernel`, `arch` matrix values
    # and will convert them to oci os/arch/variant pairs for cosign
    kernel:
    - linux
    - darwin
    arch:
    - amd64
    - armv7
  files:
  - # path to your blob file
    path: path/to/your/upload/target
    # custom content type of the blob file
    content_type: ""
  signing:
    # enable siging
    enabled: true
    # siging key
    private_key@env: ${MY_COSIGN_PRIVATE_KEY}
    # password of the siging key
    private_key_password@env: ${MY_COSIGN_PASSWORD}
    # use a different repository for signature storage (set COSIGN_REPOSITORY)
    repo: sig.example.com/dist/foo

    # verify signature after signing
    verify: true
    # user accessible public key, if not defined, will derive from the signing `key`
    public_key@http: https://example.com/cosign.pub
  # image_names is the same with docker/buildah build task's
  # omit image tag to let dukkha generate tag for you
  image_names:
  - image: example.com/dist/foo:latest-amd64
    # manifest is not supported
    # manifest: example.com/dist/foo:latest
```
