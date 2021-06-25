ARG ARCH=amd64
ARG HOST_PLATFORM_ARCH=amd64

FROM ghcr.io/arhat-dev/builder-go:alpine as builder
FROM ghcr.io/arhat-dev/go:debian-armv5
ARG APP=dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
