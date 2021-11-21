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

### Task `cosign:sign`

Sign local files

```yaml
cosign:sign:
- name: foo

  # siging key
  private_key@env: ${MY_COSIGN_PRIVATE_KEY}

  # password of the siging key
  private_key_password@env: ${MY_COSIGN_PASSWORD}

  # Verify signature just created using public key
  #
  # defaults to true
  verify: true

  # End user accessible public key, only used when verify is set to true
  #
  # if not set and `verify` is true, derive from the signing key (.private_key)
  public_key@http: https://example.com/cosign.pub

  files:
  - # path to your file
    path: path/to/local/file
    # signature output file
    # if not set, add ".sig" suffix to .path
    #output: file.sig
```

### Task `cosign:sign-image`

Sign container image already pushed to OCI registry

```yaml
cosign:sign-image:
- name: foo

  # same options in cosign:sign are not shown
  # also, there is no `files` option in cosign:sign-image

  # use a different repository for signature storage
  #
  # actually set COSIGN_REPOSITORY
  repo: sig.example.com/dist/foo

  # additional string key value pairs added when sign
  annotations:
    foo: bar

  # image_names to sign
  image_names:
  - image: example.com/dist/foo:latest-amd64
    # manifest is not supported
    # manifest: example.com/dist/foo:latest
```

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

  # signing config
  signing:
    # sign the uploaded blob
    enabled: true

    # other options are the same as signing options in cosign:sign-image (without image_names)

  # image_names for uploaded files
  #
  # omit image tag to let dukkha generate tag for you
  image_names:
  - image: example.com/dist/foo:latest-amd64
    # manifest is not supported
    # manifest: example.com/dist/foo:latest
```
