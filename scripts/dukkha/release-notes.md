# Release Notes

## Features

- Add foo support #issue-ref
- Add bar support #issue-ref

## Bug fixes

- Fixed foo #issue-ref
- Fixed bar #issue-ref #pr-ref

## Breaking Changes

- Foo ...
- Bar ...

## Changes since `{{ .Env.CHANGELOG_SINCE }}`

{{ .Env.CHANGELOG }}

## Images

- `ghcr.io/arhat-dev/dukkha:{{ .Env.GIT_TAG | trimPrefix "v" }}`

## Artifacts

Fetch signed pre-built executables using [`sget`](https://github.com/sigstore/cosign#blobs)

```bash
sget -key https://arhat.dev/.well-known/cosign.pub ghcr.io/arhat-dev/dist/dukkha:{{ .Env.GIT_TAG | trimPrefix "v" }}-{KERNEL}-{ARCH}
```
