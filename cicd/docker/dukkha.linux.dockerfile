ARG ARCH=amd64

FROM ghcr.io/arhat-dev/builder-go:alpine as builder
FROM ghcr.io/arhat-dev/go:alpine-${ARCH}
ARG APP=dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
