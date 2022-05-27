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

## Changes since `{{ env.CHANGELOG_SINCE }}`

{{ env.CHANGELOG }}

## Images

- `ghcr.io/arhat-dev/dukkha:{{ git.tag | trimPrefix "v" }}`

## Artifacts

Fetch signed pre-built executables using [`sget`](https://github.com/sigstore/cosign#blobs)

```bash
sget --key https://arhat.dev/.well-known/cosign.pub ghcr.io/arhat-dev/dist/dukkha:{{ git.tag | trimPrefix "v" }}-{KERNEL}-{ARCH}
```
