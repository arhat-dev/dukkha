# `skopeo`

__NOTE:__ This tool is currently not implemented in `dukkha`.

Manage OCI images using [`skopeo`](https://github.com/containers/skopeo)

## Supported Tasks

### Task `skopeo:copy`

```yaml
skopeo:copy:
- name: foo
  images:
  - from: example.com/foo:latest
    to: signed.example.com/foo:latest
  - from: example.com/bar:v1.0.0
    to: unsigned.example.com/bar:v1.0.0
  source:
    # remove signatures of source image, useful when the destination registry
    # doesn't support signatures
    #
    # this option takes no effect when destination.signing.enabled is true
    remove_signatures: false
    decryption:
      enabled: true
      key: ${PRIVATE_KEY}
      passphrase: ${PRIVATE_KEY_PASSPHRASE}
  destination:
    # manifest format, one of [oci, v2s1, v2s2]
    # if not set, will stay the same as source manifest
    manifest_format: oci
    signing:
      enabled: true
      pgp_key_id: ${PGP_KEY_ID}
    encryption:
      enabled: true
      key: ${PUBLIC_KEY}
      passphrase: ${PUBLIC_KEY_PASSPHRASE}
    compression:
      enabled: true
      # compression format, one of [gzip, zstd]
      format: gzip
      # compression level
      # for gzip: [1, 9]
      # for zstd: [1, 20]
      level: 9
```
