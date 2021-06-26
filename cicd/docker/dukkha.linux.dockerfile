ARG ARCH=amd64

FROM ghcr.io/arhat-dev/builder-go:alpine as builder

ARG MATRIX_ARCH
RUN set -ex ;\
    make dukkha && \
    ./build/dukkha golang build dukkha -m kernel=linux -m arch=${MATRIX_ARCH}

FROM ghcr.io/arhat-dev/go:alpine-${ARCH}

ARG MATRIX_ARCH
COPY --from=builder /app/build/dukkha.linux.${MATRIX_ARCH} /dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
