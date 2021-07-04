ARG MATRIX_ARCH
ARG MATRIX_ROOTFS

FROM ghcr.io/arhat-dev/builder-go:alpine as builder

ARG MATRIX_ARCH

COPY . /app
RUN set -ex ;\
    make dukkha && \
    ./build/dukkha golang build dukkha -m kernel=linux -m arch=${MATRIX_ARCH}

FROM ghcr.io/arhat-dev/go:${MATRIX_ROOTFS}-${MATRIX_ARCH}

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ARG MATRIX_ARCH
COPY --from=builder /app/build/dukkha.linux.${MATRIX_ARCH} /dukkha

ENTRYPOINT [ "/dukkha" ]
