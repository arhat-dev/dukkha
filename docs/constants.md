# Constants

Wellknown values in `dukkha` to ease cross-platform task execution

## System Kernel

- `windows`
- `linux`
- `darwin`
- `freebsd`
- `netbsd`
- `openbsd`
- `solaris`
- `illumos`
- `js`
- `aix`
- `android`
- `ios`
- `plan9`

__NOTE:__ These values are the same as golang `GOOS` values

## System Arch

- `x86`
- `amd64`
- `amd64v1` (alias of amd64)
- `amd64v2` (2009+)
- `amd64v3` (2015+)
- `amd64v4` (avx512 extension)
- `arm64`
- `armv5`
- `armv6`
- `armv7`
- `mips`
- `mipssf`
- `mipsle`
- `mipslesf`
- `mips64`
- `mips64sf`
- `mips64le`
- `mips64lesf`
- `ppc`
- `ppcsf`
- `ppcle`
- `ppclesf`
- `ppc64`
- `ppc64le`
- `ppc64v8`
- `ppc64v8le`
- `ppc64v9`
- `ppc64v9le`
- `riscv64`
- `s390x`
- `ia64`

__NOTE:__ These values are defined in package [arhat.dev/pkg/archconst][https://github.com/arhat-dev/go-pkg/blob/master/archconst/values.go]

- `sf` means softfloat, defaults to hardfloat if available in that cpu arch
- `le` means little endian, defaults to arch default (e.g. `arm64` implies little endian, but `ppc64` implies big endian)

## Libc

- `gnu`
- `musl`
- `msvc`
