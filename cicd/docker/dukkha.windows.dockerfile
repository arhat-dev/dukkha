ARG MATRIX_ARCH

FROM ghcr.io/arhat-dev/builder-golang:1.16-alpine as builder

ARG MATRIX_ARCH

COPY . /app
RUN set -ex ;\
    make dukkha && \
    ./build/dukkha golang local build dukkha -m kernel=linux -m arch=${MATRIX_ARCH}

# TODO: support multiarch build
FROM mcr.microsoft.com/windows/servercore:ltsc2019

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ARG MATRIX_ARCH
COPY --from=builder /app/build/dukkha.linux.${MATRIX_ARCH} /dukkha

ENTRYPOINT [ "/dukkha" ]
