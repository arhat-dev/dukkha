ARG ARCH=amd64
ARG HOST_PLATFORM_ARCH=amd64

FROM --platform=linux/${HOST_PLATFORM_ARCH} ghcr.io/arhat-dev/builder-go:alpine as builder
FROM --platform=linux/mips64le ghcr.io/arhat-dev/go:debian-mips64le
ARG APP=dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
