ARG ARCH=amd64
ARG HOST_PLATFORM_ARCH=amd64
ARG PLATFORM_ARCH=amd64

FROM --platform=linux/${HOST_PLATFORM_ARCH} ghcr.io/arhat-dev/builder-go:alpine as builder
FROM --platform=linux/${PLATFORM_ARCH} ghcr.io/arhat-dev/go:alpine-${ARCH}
ARG APP=dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
